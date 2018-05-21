package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smart-cs/scheduler-backend/models"
)

func TestCombinationsNoConflict(t *testing.T) {
	assert := assert.New(t)
	ch := models.CourseHelper{}

	combinations := [][]models.CourseSection{
		[]models.CourseSection{
			{
				Name: "MATH 100 101",
				Sessions: []models.ClassSession{
					{
						Activity: "Lecture",
						Term:     "1",
						Day:      "Mon Wed Fri",
						Start:    800,
						End:      900,
					},
				},
			},
		},
	}
	combinations = ch.CombinationsNoConflict(combinations, []models.CourseSection{
		{
			Name: "MATH 101 201",
			Sessions: []models.ClassSession{
				{
					Activity: "Lecture",
					Term:     "2",
					Day:      "Mon Wed Fri",
					Start:    800,
					End:      900,
				},
			},
		},
	})
	assert.Len(combinations, 1)
	combinations = ch.CombinationsNoConflict(combinations, []models.CourseSection{
		{
			Name: "CPSC 100 101",
			Sessions: []models.ClassSession{
				{
					Activity: "Lecture",
					Term:     "1",
					Day:      "Mon Wed Fri",
					Start:    800,
					End:      900,
				},
			},
		},
	})
	assert.Len(combinations, 0)
}

func TestIsIncluded(t *testing.T) {
	assert := assert.New(t)
	ch := models.CourseHelper{}

	assert.True(ch.IsIncluded("Lecture", []models.ActivityType{models.Lecture}))
	assert.True(ch.IsIncluded("Lecture", []models.ActivityType{models.Lecture, models.Laboratory}))
	assert.False(ch.IsIncluded("Lecture", []models.ActivityType{models.Laboratory}))
	assert.False(ch.IsIncluded("Lecture", []models.ActivityType{}))
}

func TestConflictInSchedule(t *testing.T) {
	assert := assert.New(t)
	ch := models.CourseHelper{}

	schedule := models.Schedule{
		Courses: []models.CourseSection{
			{
				Name: "MATH 100 101",
				Sessions: []models.ClassSession{
					{
						Activity: "Lecture",
						Term:     "1",
						Day:      "Mon Wed Fri",
						Start:    800,
						End:      900,
					},
				},
			},
			{
				Name: "MATH 101 201",
				Sessions: []models.ClassSession{
					{
						Activity: "Lecture",
						Term:     "2",
						Day:      "Mon Wed Fri",
						Start:    800,
						End:      900,
					},
				},
			},
		},
	}

	assert.False(ch.ConflictInSchedule(models.Schedule{}))
	assert.False(ch.ConflictInSchedule(schedule))
	schedule.Courses = append(schedule.Courses, models.CourseSection{
		Name: "CPSC 100 101",
		Sessions: []models.ClassSession{
			{
				Activity: "Lecture",
				Term:     "1",
				Day:      "Mon Wed Fri",
				Start:    800,
				End:      900,
			},
		},
	})
	assert.True(ch.ConflictInSchedule(schedule))
}
