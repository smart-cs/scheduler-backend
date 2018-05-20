package schedules

import (
	"github.com/smart-cs/scheduler-backend/database"
	"github.com/smart-cs/scheduler-backend/models"
)

// ScheduleCreator is the interface to create schedules.
type ScheduleCreator interface {
	Create(courses []string, options ScheduleSelectOptions) []models.Schedule
}

// DefaultScheduleCreator implements ScheduleCreator.
type DefaultScheduleCreator struct {
	ds     database.Datastore
	helper models.CourseHelper
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
		ds:     database.NewDatastore(),
		helper: models.CourseHelper{},
	}
}

// Create returns all non-conflicting schedules given a list of courses.
func (sc *DefaultScheduleCreator) Create(courses []string, options ScheduleSelectOptions) []models.Schedule {
	var schedules []models.Schedule
	for _, c := range courses {
		// Skip invalid courses.
		if !sc.ds.CourseExists(c) {
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

func (sc *DefaultScheduleCreator) addCourseToSchedules(schedules []models.Schedule, c, term string, selectLabsAndTuts bool) ([]models.Schedule, bool) {
	lectureSections := sc.ds.GetSections(c, term, models.Lecture, models.Seminar, models.Studio)
	hasLabs := sc.ds.CourseHasSectionWithActivity(c, models.Laboratory)
	hasTuts := sc.ds.CourseHasSectionWithActivity(c, models.Tutorial)

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
		labSections := sc.ds.GetSections(c, term, models.Laboratory)
		sectionsArray = sc.helper.CombinationsNoConflict(sectionsArray, labSections)
	}
	if hasTuts {
		tutSections := sc.ds.GetSections(c, term, models.Tutorial)
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
