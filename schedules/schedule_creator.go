package schedules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/smart-cs/scheduler-backend/database"
	"github.com/smart-cs/scheduler-backend/models"
)

// ScheduleCreator is the interface to create schedules.
type ScheduleCreator interface {
	Create(courses []string, options ScheduleSelectOptions) []models.Schedule
}

// DefaultScheduleCreator implements ScheduleCreator.
type DefaultScheduleCreator struct {
	db     database.CourseDatabase
	helper CourseHelper
}

// ScheduleSelectOptions is a criteria for selecting schedules.
type ScheduleSelectOptions struct {
	// Term must be 1, 2, 1-2
	Term                   string
	SelectLabsAndTutorials bool
}

// NewScheduleCreator constructs a new ScheduleCreator.
func NewScheduleCreator() ScheduleCreator {
	return &DefaultScheduleCreator{
		db:     database.CourseDB(),
		helper: CourseHelper{},
	}
}

// Create returns all non-conflicting schedules given a list of courses.
func (sc *DefaultScheduleCreator) Create(courses []string, options ScheduleSelectOptions) []models.Schedule {
	var schedules []models.Schedule
	for _, c := range courses {
		// Skip invalid courses.
		if !sc.courseExists(c) {
			continue
		}

		if options.Term == "1-2" {
			newSchedulesTerm1, addedTerm1 := sc.addCourseToSchedules(schedules, c, "1", options.SelectLabsAndTutorials)
			newSchedulesTerm2, addedTerm2 := sc.addCourseToSchedules(schedules, c, "2", options.SelectLabsAndTutorials)
			if !addedTerm1 && !addedTerm2 {
				return []models.Schedule{}
			}

			schedules = append(newSchedulesTerm1, newSchedulesTerm2...)
			continue
		}

		newSchedules, added := sc.addCourseToSchedules(schedules, c, options.Term, options.SelectLabsAndTutorials)
		if !added {
			return []models.Schedule{}
		}
		schedules = newSchedules
	}
	return schedules
}

func (sc *DefaultScheduleCreator) courseExists(course string) bool {
	dept := strings.Split(course, " ")[0]
	_, present := sc.db[dept][course]
	return present
}

func (sc *DefaultScheduleCreator) courseHasActivity(course string, activity models.ActivityType) bool {
	dept := strings.Split(course, " ")[0]
	for _, section := range sc.db[dept][course] {
		if len(section.Activity) == 0 {
			// TODO: Handle invalid input.
			continue
		}

		if section.Activity[0] == activity.String() {
			return true
		}
	}
	return false
}

func (sc *DefaultScheduleCreator) addCourseToSchedules(schedules []models.Schedule, c, term string, selectLabsAndTuts bool) ([]models.Schedule, bool) {
	lectureSections := sc.createSections(c, term, models.Lecture, models.Seminar, models.Studio)
	hasLabs := sc.courseHasActivity(c, models.Laboratory)
	hasTuts := sc.courseHasActivity(c, models.Tutorial)

	if !selectLabsAndTuts || (!hasLabs && !hasTuts) {
		// Just add the lecture sections.
		schedules = sc.addSections(schedules, lectureSections)
		return schedules, len(schedules) != 0
	}

	// Add sections including Labs and Tutorials.
	var sectionsArray [][]models.CourseSection
	for _, section := range lectureSections {
		sectionsArray = append(sectionsArray, []models.CourseSection{section})
	}
	if hasLabs {
		labSections := sc.createSections(c, term, models.Laboratory)
		sectionsArray = sc.helper.CombinationsNoConflict(sectionsArray, labSections)
	}
	if hasTuts {
		tutSections := sc.createSections(c, term, models.Tutorial)
		sectionsArray = sc.helper.CombinationsNoConflict(sectionsArray, tutSections)
	}
	schedules = sc.addSectionBlocks(schedules, sectionsArray)
	return schedules, len(schedules) != 0
}

func (sc *DefaultScheduleCreator) addSectionBlocks(schedules []models.Schedule, sectionsArray [][]models.CourseSection) []models.Schedule {
	if len(schedules) == 0 {
		for _, sections := range sectionsArray {
			schedules = append(schedules, models.Schedule{Courses: sections})
		}
		return schedules
	}

	newSchedules := []models.Schedule{}
	addedASection := false
	for _, schedule := range schedules {
		for _, sections := range sectionsArray {
			newSchedule, added := sc.addSection(schedule, sections...)
			if added {
				newSchedules = append(newSchedules, newSchedule)
				addedASection = true
			}
		}
	}
	if !addedASection {
		return []models.Schedule{}
	}
	return newSchedules
}

func (sc *DefaultScheduleCreator) addSections(schedules []models.Schedule, sections []models.CourseSection) []models.Schedule {
	if len(schedules) == 0 {
		for _, section := range sections {
			sections := []models.CourseSection{section}
			schedules = append(schedules, models.Schedule{
				Courses: sections,
			})
		}
		return schedules
	}

	newSchedules := []models.Schedule{}
	for _, schedule := range schedules {
		for _, section := range sections {
			newSchedule, added := sc.addSection(schedule, section)
			if added {
				newSchedules = append(newSchedules, newSchedule)
			}
		}
	}
	return newSchedules
}

// addSection returns the new schedule if all sections can be added, otherwise returns the old schedule.
func (sc *DefaultScheduleCreator) addSection(schedule models.Schedule, sections ...models.CourseSection) (models.Schedule, bool) {
	newSchedule := schedule
	for _, section := range sections {
		newSchedule.Courses = append(newSchedule.Courses, section)
		if sc.helper.ConflictInSchedule(newSchedule) {
			return schedule, false
		}
	}
	return newSchedule, true
}

// createSections returns sections of a course with one of the specified types, thats in terms.
// Possible terms: 1, 2, 1-2.
func (sc *DefaultScheduleCreator) createSections(course, term string, activityTypes ...models.ActivityType) []models.CourseSection {
	// Course format i.e. CPSC 121
	var sections []models.CourseSection
	dept := strings.Split(course, " ")[0]
	// Go through all sections for this course.
	for sectionName, s := range sc.db[dept][course] {
		if !sc.helper.IsIncluded(s.Activity[0], activityTypes) {
			continue
		}

		if (term == "1" || term == "2") && s.Term[0] != term {
			continue
		}

		sessions, err := sc.sessions(s)
		if err != nil {
			fmt.Printf("ERROR: validating fields for %q: %s\n", sectionName, err.Error())
		}

		section := models.CourseSection{
			Name:     sectionName,
			Sessions: sessions,
		}
		sections = append(sections, section)
	}
	return sections
}

func (sc *DefaultScheduleCreator) sessions(s database.Section) ([]models.ClassSession, error) {
	var sessions []models.ClassSession
	for i, dayStr := range s.Days {
		// dayStr looks like "Mon Wed Fri".
		for _, day := range strings.Split(dayStr, " ") {
			start, err := sc.parseTime(s.StartTime[0])
			if err != nil {
				return nil, errors.Wrap(err, "no startTime")
			}
			end, err := sc.parseTime(s.EndTime[0])
			if err != nil {
				return nil, errors.Wrap(err, "no endTime")
			}
			session := models.ClassSession{
				Activity: s.Activity[i],
				Term:     s.Term[i],
				Day:      day,
				Start:    start,
				End:      end,
			}
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

// praseTime parses time in the format HH:MM to an int HHMM.
func (sc *DefaultScheduleCreator) parseTime(time string) (int, error) {
	parsed, err := strconv.Atoi(strings.Replace(time, ":", "", -1))
	return parsed, err
}
