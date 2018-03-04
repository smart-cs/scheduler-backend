package server_test

import (
	"testing"

	"github.com/nickwu241/schedulecreator-backend/server"
	"github.com/stretchr/testify/assert"
)

func TestScheduleCreator_Create(t *testing.T) {
	assert := assert.New(t)
	table := []struct {
		inCourses []string
		outLen    int
	}{
		{[]string{"APSC 201"}, 7},
		{[]string{"BIOL 111"}, 2},
		{[]string{"MATH 220"}, 9},
		{[]string{"MATH 253"}, 6},
		{[]string{"MATH 335"}, 2},
		{[]string{"APBI 260"}, 1},
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
	for _, item := range table {
		schedules := sc.Create(item.inCourses)
		assert.Equalf(
			item.outLen,
			len(schedules),
			"creating schedules from %v should return %d schedules, but got %d",
			item.inCourses,
			item.outLen,
			len(schedules),
		)
	}
}
