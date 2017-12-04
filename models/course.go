package models

// ClassSession holds time information about about a single class.
type ClassSession struct {
	// Acitvity is the type of class. e.g. 'Lecture'
	Activity string `json:"activity"`
	// Term '1' or '2' or '1-2'.
	Term string `json:"term"`
	// Day of the week. e.g. 'Mon' 'Tue' 'Wed'
	Day string `json:"day"`
	// Start time of class (24 hour representation). e.g. 1230
	Start int `json:"start"`
	// End time of class (24 hour representation). e.g. 1530
	End int `json:"end"`
}

// CourseSection represents a course section.
type CourseSection struct {
	// Name of the course: <DEPARTMENT> <LEVEL> <SECTION>. e.g. 'CPSC 121 101'
	Name string `json:"name"`
	// List of ClassSession.
	Sessions []ClassSession `json:"sessions"`
}

// Schedule represents a schedule of courses.
type Schedule struct {
	// List of Course.
	Courses []CourseSection `json:"courses"`
}
