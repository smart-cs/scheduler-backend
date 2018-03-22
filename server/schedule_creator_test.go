package server_test

import (
	"testing"

	"github.com/smart-cs/scheduler-backend/server"
	"github.com/stretchr/testify/assert"
)

func TestScheduleCreator_Create(t *testing.T) {
	assert := assert.New(t)
	testTables := []struct {
		courses   []string
		expOutLen int
	}{
		{[]string{"APSC 201"}, 7},
		{[]string{"BIOL 111"}, 2},
		{[]string{"MATH 220"}, 9},
		{[]string{"MATH 253"}, 6},
		{[]string{"MATH 335"}, 2},
		{[]string{"APBI 260"}, 1},
		{[]string{"CPEN 221"}, 1},
		{[]string{"ASIA 100"}, 2},
		{[]string{"APBI 260", "ASIA 100"}, 1},
		{[]string{"MATH 220", "MATH 253"}, 54},
		{[]string{"MATH 220", "MATH 335"}, 18},
		{[]string{"CPSC 110"}, 8},
		{[]string{"APSC 210"}, 0},
		{[]string{"non-existent-course 101"}, 0},
		{[]string{"BIOL 111", "non-existent-course 101"}, 2},
		{[]string{"non-existent-course 101", "BIOL 111"}, 2},
	}

	server.LoadLocalDatabase("coursedb.json")
	sc := server.NewScheduleCreator()
	for _, tt := range testTables {
		schedules := sc.Create(tt.courses, server.ScheduleSelectOptions{Term: "1-2"})
		assert.Equalf(
			tt.expOutLen,
			len(schedules),
			"creating schedules from %v should return %d schedules, but got %d",
			tt.courses,
			tt.expOutLen,
			len(schedules),
		)
	}
}

func TestShceduleCreator_Term(t *testing.T) {
	assert := assert.New(t)
	testTables := []struct {
		courses   []string
		term      string
		expOutLen int
	}{
		{[]string{"CPEN 221"}, "1", 1},
		{[]string{"CPEN 221"}, "2", 0},
	}

	server.LoadLocalDatabase("coursedb.json")
	sc := server.NewScheduleCreator()
	for _, tt := range testTables {
		schedules := sc.Create(tt.courses, server.ScheduleSelectOptions{Term: tt.term})
		assert.Equalf(
			tt.expOutLen,
			len(schedules),
			"creating schedules from %v should return %d schedules, but got %d",
			tt.courses,
			tt.expOutLen,
			len(schedules),
		)
	}
}
