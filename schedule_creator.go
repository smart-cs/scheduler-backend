package main

import (
	"strconv"
	"strings"

	"github.com/nickwu241/schedulecreator-backend/models"
)

// ScheduleCreator is the interface to create schedules.
type ScheduleCreator interface {
	Create(courses []string) []models.Schedule
}

// DefaultScheduleCreator implments ScheduleCreator.
type DefaultScheduleCreator struct {
	db CourseDatabase
}

// NewScheduleCreator constructs a new ScheduleCreator.
func NewScheduleCreator() ScheduleCreator {
	return &DefaultScheduleCreator{
		db: CourseDB(),
	}
}

// Create returns all non-conflicting schedules given a list of courses.
func (d *DefaultScheduleCreator) Create(courses []string) []models.Schedule {
	var schedules []models.Schedule
	for _, c := range courses {
		// Skip invalid courses.
		if !d.courseExists(c) {
			continue
		}
		schedules = d.addSections(schedules, d.createSections(c))
	}
	return schedules
}

func (d *DefaultScheduleCreator) courseExists(course string) bool {
	dept := strings.Split(course, " ")[0]
	_, present := d.db[dept][course]
	return present
}

func (d *DefaultScheduleCreator) addSections(schedules []models.Schedule, sections []models.CourseSection) []models.Schedule {
	var result []models.Schedule
	if len(schedules) == 0 {
		for _, section := range sections {
			sections := []models.CourseSection{section}
			result = append(result, models.Schedule{
				Courses: sections,
			})
		}
		return result
	}

	for _, schedule := range schedules {
		for _, section := range sections {
			new := schedule
			new.Courses = append(new.Courses, section)
			// Only add the new course if it doesn't conflict.
			if d.conflictInSchedule(new) {
				continue
			}
			result = append(result, new)
		}
	}
	return result
}

func (d *DefaultScheduleCreator) createSections(course string) []models.CourseSection {
	var sections []models.CourseSection
	dept := strings.Split(course, " ")[0]
	// Go through all sections for this course.
	for sectionName, s := range d.db[dept][course] {
		// TODO: do we want to check for other activities?
		if s.Activity[0] != "Lecture" &&
			s.Activity[0] != "Seminar" &&
			s.Activity[0] != "Studio" {
			continue
		}
		// Create the sessions for each section.
		var sessions []models.ClassSession
		for i, dayStr := range s.Days {
			// dayStr looks like "Mon Wed Fri".
			for _, day := range strings.Split(dayStr, " ") {

				// TODO: refactor this logic out.
				start, err := strconv.Atoi(strings.Replace(s.StartTime[i], ":", "", -1))
				if err != nil {
					panic(err)
				}
				end, err := strconv.Atoi(strings.Replace(s.EndTime[i], ":", "", -1))
				if err != nil {
					panic(err)
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
		section := models.CourseSection{
			Name:     sectionName,
			Sessions: sessions,
		}
		sections = append(sections, section)
	}
	return sections
}

func (d *DefaultScheduleCreator) conflictSession(s1 models.ClassSession, s2 models.ClassSession) bool {
	return s1.Term == s2.Term && s1.Day == s2.Day &&
		((s1.Start <= s2.Start && s2.Start < s1.End) ||
			(s1.Start < s2.End && s2.End <= s1.End))
}

func (d *DefaultScheduleCreator) conflictSection(s1 models.CourseSection, s2 models.CourseSection) bool {
	for _, ses1 := range s1.Sessions {
		for _, ses2 := range s2.Sessions {
			if d.conflictSession(ses1, ses2) {
				return true
			}
		}
	}
	return false
}

func (d *DefaultScheduleCreator) conflictInSchedule(schedule models.Schedule) bool {
	for _, c1 := range schedule.Courses {
		for _, c2 := range schedule.Courses {
			if c1.Name != c2.Name && d.conflictSection(c1, c2) {
				return true
			}
		}
	}
	return false
}
