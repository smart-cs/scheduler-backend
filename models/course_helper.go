package models

// CourseHelper contains helpful operations on course models.
type CourseHelper struct{}

// CombinationsNoConflict generates all the combinations of CourseSections that doesn't conflict.
func (c *CourseHelper) CombinationsNoConflict(result [][]CourseSection, sections []CourseSection) [][]CourseSection {
	var newResult [][]CourseSection
	for _, comb := range result {
		for _, section := range sections {
			// Create an array to use conflictInSections.
			if c.conflictInSections(comb, []CourseSection{section}) {
				continue
			}
			newComb := append(comb, section)
			newResult = append(newResult, newComb)
		}
	}
	return newResult
}

// IsIncluded returns true if desiredTypes contains the activity.
func (c *CourseHelper) IsIncluded(activity string, desiredTypes []ActivityType) bool {
	for _, a := range desiredTypes {
		if activity == a.String() {
			return true
		}
	}
	return false
}

// ConflictInSchedule returns true if there is a conflict in the schedule.
func (c *CourseHelper) ConflictInSchedule(schedule Schedule) bool {
	return c.conflictInSections(schedule.Courses, schedule.Courses)
}

func (c *CourseHelper) conflictInSections(s1s, s2s []CourseSection) bool {
	for _, s1 := range s1s {
		for _, s2 := range s2s {
			if s1.Name != s2.Name && c.conflictSection(s1, s2) {
				return true
			}
		}
	}
	return false
}

func (c *CourseHelper) conflictSection(s1, s2 CourseSection) bool {
	for _, ses1 := range s1.Sessions {
		for _, ses2 := range s2.Sessions {
			if c.conflictSession(ses1, ses2) {
				return true
			}
		}
	}
	return false
}

func (c *CourseHelper) conflictSession(s1, s2 ClassSession) bool {
	return s1.Term == s2.Term && s1.Day == s2.Day &&
		((s1.Start <= s2.Start && s2.Start < s1.End) ||
			(s1.Start < s2.End && s2.End <= s1.End))
}
