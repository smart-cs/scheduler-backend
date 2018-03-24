package server_test

import (
	"testing"

	"github.com/smart-cs/scheduler-backend/server"
	"github.com/stretchr/testify/assert"
)

func setup() {
	server.LoadLocalDatabase("coursedb.json")
}

type scheduleCreatorTestTable struct {
	courses         []string
	expSchedulesLen int
	expCoursesLen   int
}

var defaultTestTables = []scheduleCreatorTestTable{
	// Regular cases.
	{[]string{"APSC 201"}, 7, 1},
	{[]string{"ASIA 100"}, 2, 1},
	{[]string{"BIOL 111"}, 2, 1},
	{[]string{"MATH 220"}, 9, 1},
	{[]string{"MATH 253"}, 6, 1},
	{[]string{"MATH 335"}, 2, 1},
	{[]string{"MATH 220", "MATH 253"}, 54, 2},
	{[]string{"MATH 220", "MATH 335"}, 18, 2},
	// Special cases below.
	{[]string{"APSC 210"}, 0, 0}, // Co-op placement.
	{[]string{"non-existent-course 101"}, 0, 0},
	{[]string{"BIOL 111", "non-existent-course 101"}, 2, 1},
	{[]string{"non-existent-course 101", "BIOL 111"}, 2, 1},
}

func assertTables(assert *assert.Assertions, testTables []scheduleCreatorTestTable, options server.ScheduleSelectOptions) {
	sc := server.NewScheduleCreator()
	for _, tt := range testTables {
		schedules := sc.Create(tt.courses, options)
		assert.Equalf(
			tt.expSchedulesLen, len(schedules),
			"creating schedules from %v should return %d schedules, but got %d",
			tt.courses, tt.expSchedulesLen, len(schedules),
		)
		for _, schedule := range schedules {
			if !assert.Equalf(
				tt.expCoursesLen, len(schedule.Courses),
				"schedule %v should contain %d courses, but got %d",
				schedule, tt.expCoursesLen, len(schedule.Courses)) {
				break
			}
		}
	}
}

func TestScheduleCreator_CreateDefault(t *testing.T) {
	setup()
	testTables := append(defaultTestTables, []scheduleCreatorTestTable{
		{[]string{"APBI 260"}, 1, 1},
		{[]string{"CPEN 221"}, 1, 1},
		{[]string{"CPSC 110"}, 8, 1},
		{[]string{"CPSC 210"}, 6, 1},
		{[]string{"CPSC 221"}, 5, 1},
		{[]string{"APBI 260", "ASIA 100"}, 1, 2},
	}...)
	options := server.ScheduleSelectOptions{
		Term: "1-2",
		SelectLabsAndTutorials: false,
	}
	assertTables(assert.New(t), testTables, options)
}

func TestScheduleCreator_CreateWithLabsAndTuts(t *testing.T) {
	setup()
	testTables := append(defaultTestTables, []scheduleCreatorTestTable{
		{[]string{"APBI 260"}, 1, 2},
		{[]string{"CPEN 221"}, 5, 3},
		{[]string{"CPSC 110"}, 81, 3},
		{[]string{"CPSC 210"}, 99, 2},
		{[]string{"CPSC 221"}, 72, 2},
		{[]string{"APBI 260", "ASIA 100"}, 0, 0},
	}...)
	options := server.ScheduleSelectOptions{
		Term: "1-2",
		SelectLabsAndTutorials: true,
	}
	assertTables(assert.New(t), testTables, options)
}

func TestShceduleCreator_CreateWithTerm(t *testing.T) {
	setup()
	assert := assert.New(t)
	testTables := []struct {
		courses         []string
		term            string
		expSchedulesLen int
	}{
		{[]string{"CPEN 221"}, "1", 1},
		{[]string{"CPEN 221"}, "2", 0},
	}

	sc := server.NewScheduleCreator()
	for _, tt := range testTables {
		schedules := sc.Create(tt.courses, server.ScheduleSelectOptions{Term: tt.term})
		assert.Equalf(
			tt.expSchedulesLen, len(schedules),
			"creating schedules from %v should return %d schedules, but got %d",
			tt.courses, tt.expSchedulesLen, len(schedules),
		)
	}
}
