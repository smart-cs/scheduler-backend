package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/smart-cs/scheduler-backend/models"
)

// ScheduleCreator is the interface to create schedules.
type ScheduleCreator interface {
	Create(courses []string) []models.Schedule
}

// DefaultScheduleCreator implements ScheduleCreator.
type DefaultScheduleCreator struct {
	db     CourseDatabase
	helper CourseHelper
}

// NewScheduleCreator constructs a new ScheduleCreator.
func NewScheduleCreator() ScheduleCreator {
	return &DefaultScheduleCreator{
		db:     CourseDB(),
		helper: CourseHelper{},
	}
}

// Create returns all non-conflicting schedules given a list of courses.
func (sc *DefaultScheduleCreator) Create(courses []string) []models.Schedule {
	var schedules []models.Schedule
	for _, c := range courses {
		// Skip invalid courses.
		if !sc.courseExists(c) {
			continue
		}
		lectureTypes := []models.ActivityType{models.Lecture, models.Seminar, models.Studio}
		schedules = sc.addSections(schedules, sc.createSections(c, lectureTypes))
		// schedules = d.addSections(schedules, d.createSections(c, []models.ActivityType{models.Laboratory}))
		// schedules = d.addSections(schedules, d.createSections(c, []models.ActivityType{models.Tutorial}))
	}
	return schedules
}

func (sc *DefaultScheduleCreator) courseExists(course string) bool {
	dept := strings.Split(course, " ")[0]
	_, present := sc.db[dept][course]
	return present
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

// addSection returns the new schedule if it succeeds, old schedule if it fails
func (sc *DefaultScheduleCreator) addSection(schedule models.Schedule, section models.CourseSection) (models.Schedule, bool) {
	newSchedule := schedule
	newSchedule.Courses = append(newSchedule.Courses, section)
	if sc.helper.ConflictInSchedule(newSchedule) {
		return schedule, false
	}
	return newSchedule, true
}

func (sc *DefaultScheduleCreator) createSections(course string, activityTypes []models.ActivityType) []models.CourseSection {
	// Course format i.e. CPSC 121
	var sections []models.CourseSection
	dept := strings.Split(course, " ")[0]
	// Go through all sections for this course.
	for sectionName, s := range sc.db[dept][course] {
		if !sc.helper.IsIncluded(s.Activity[0], activityTypes) {
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

func (sc *DefaultScheduleCreator) sessions(s Section) ([]models.ClassSession, error) {
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
