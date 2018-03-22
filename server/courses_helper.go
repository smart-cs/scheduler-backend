package server

import "github.com/smart-cs/scheduler-backend/models"

// CourseHelper contains helpful operations on course models.
type CourseHelper struct{}

// IsIncluded returns true if desiredTypes contains the activity.
func (c *CourseHelper) IsIncluded(activity string, desiredTypes []models.ActivityType) bool {
	for _, a := range desiredTypes {
		if activity == a.String() {
			return true
		}
	}
	return false
}

// ConflictInSchedule returns true if there is a conflict in the schedule.
func (c *CourseHelper) ConflictInSchedule(schedule models.Schedule) bool {
	for _, c1 := range schedule.Courses {
		for _, c2 := range schedule.Courses {
			if c1.Name != c2.Name && c.conflictSection(c1, c2) {
				return true
			}
		}
	}
	return false
}

func (c *CourseHelper) conflictSection(s1 models.CourseSection, s2 models.CourseSection) bool {
	for _, ses1 := range s1.Sessions {
		for _, ses2 := range s2.Sessions {
			if c.conflictSession(ses1, ses2) {
				return true
			}
		}
	}
	return false
}

func (c *CourseHelper) conflictSession(s1 models.ClassSession, s2 models.ClassSession) bool {
	return s1.Term == s2.Term && s1.Day == s2.Day &&
		((s1.Start <= s2.Start && s2.Start < s1.End) ||
			(s1.Start < s2.End && s2.End <= s1.End))
}
