package models_test

import (
	"testing"

	"github.com/nickwu241/schedulecreator-backend/models"

	"github.com/stretchr/testify/assert"
)

func TestActivityType(t *testing.T) {
	assert := assert.New(t)
	table := []struct {
		in  models.ActivityType
		out string
	}{
		{models.Lecture, "Lecture"},
		{models.Laboratory, "Laboratory"},
		{models.Seminar, "Seminar"},
		{models.Studio, "Studio"},
		{models.Tutorial, "Tutorial"},
	}

	for _, item := range table {
		assert.Equalf(
			item.out,
			item.in.String(),
			"activity type %v should be equal to %q",
			item.in.String(),
			item.out,
		)
	}
}
