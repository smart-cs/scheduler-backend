package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/nickwu241/schedulecreator-backend/models"
)

const dbFilename = "coursedb.json"

// CourseHelper provides utility functions about schedules and courses.
type CourseHelper interface {
	ValidDepartments() []string
	ValidCourses() []string
	CreateSchedules(courses []string) []models.Schedule
}

// DefaultCourseHelper is the default implementation of CourseHelper.
type DefaultCourseHelper struct{}

var courseDB = initDatabase()
var validDepartments []string
var validCourses []string

// CourseDatabase is the schema for our courses database.
// Schema: DEPARTMENT_NAME -> COURSE_NAME -> COURSE_SECTION_NAME -> Section
// e.g. DB["CPSC"]["CPSC 121"]["CPSC 121 101"] to get the underyling Section
type CourseDatabase map[string]map[string]map[string]Section

// Section is a Section of a UBC course.
type Section struct {
	Activity  []string `json:"activity"`
	Days      []string `json:"days"`
	EndTime   []string `json:"end_time"`
	Interval  string   `json:"interval"`
	StartTime []string `json:"start_time"`
	Status    string   `json:"status"`
	Term      []string `json:"term"`
}

// NewCourseHelper constructs a new CourseHelper.
func NewCourseHelper() CourseHelper {
	for deptarment, sectionMap := range courseDB {
		validDepartments = append(validDepartments, deptarment)
		for sectionName := range sectionMap {
			validCourses = append(validCourses, sectionName)
		}
	}
	return &DefaultCourseHelper{}
}

// ValidDepartments returns the valid departments.
func (d *DefaultCourseHelper) ValidDepartments() []string {
	return validDepartments
}

// ValidCourses returns the valid courses.
func (d *DefaultCourseHelper) ValidCourses() []string {
	return validCourses
}

// CreateSchedules returns all non-conflicting schedules given a list of courses
func (d *DefaultCourseHelper) CreateSchedules(courses []string) []models.Schedule {
	var schedules []models.Schedule
	for _, c := range courses {
		schedules = addSections(schedules, createSections(c))
	}
	return schedules
}

func initDatabase() CourseDatabase {
	_, callerPath, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller information")
	}

	dbPath := path.Join(path.Dir(callerPath), dbFilename)
	f, err := os.Open(dbPath)
	if err != nil {
		panic("can't initialize database")
	}

	var db CourseDatabase
	json.NewDecoder(bufio.NewReader(f)).Decode(&db)
	return db
}

func addSections(schedules []models.Schedule, sections []models.CourseSection) []models.Schedule {
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
			// Only add the new course if it doesn't conflict
			if conflictInSchedule(new) {
				continue
			}
			result = append(result, new)
		}
	}
	return result
}

func createSections(course string) []models.CourseSection {
	var sections []models.CourseSection
	dept := strings.Split(course, " ")[0]
	// Go through all sections for this course.
	for sectionName, s := range courseDB[dept][course] {
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

				// TODO: refacotr this logic out
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

func conflictSession(s1 models.ClassSession, s2 models.ClassSession) bool {
	return s1.Term == s2.Term && s1.Day == s2.Day &&
		((s1.Start <= s2.Start && s2.Start < s1.End) ||
			(s1.Start < s2.End && s2.End <= s1.End))
}

func conflictSection(s1 models.CourseSection, s2 models.CourseSection) bool {
	for _, ses1 := range s1.Sessions {
		for _, ses2 := range s2.Sessions {
			if conflictSession(ses1, ses2) {
				return true
			}
		}
	}
	return false
}

func conflictInSchedule(schedule models.Schedule) bool {
	for _, c1 := range schedule.Courses {
		for _, c2 := range schedule.Courses {
			if c1.Name != c2.Name && conflictSection(c1, c2) {
				return true
			}
		}
	}
	return false
}
