package schedules

import (
	"strings"

	"github.com/derekparker/trie"
	"github.com/smart-cs/scheduler-backend/database"
)

// AutoCompleter finds courses with certain prefixes.
type AutoCompleter interface {
	CoursesWithPrefix(prefix string) []string
}

// DefaultAutoCompleter implements AutoCompleter.
type DefaultAutoCompleter struct {
	Courses trie.Trie
}

// NewAutoCompleter constructs an AutoCompleter.
func NewAutoCompleter() AutoCompleter {
	t := trie.New()
	for _, d := range database.ValidCourses() {
		t.Add(d, nil)
	}
	return &DefaultAutoCompleter{
		Courses: *t,
	}
}

// CoursesWithPrefix finds all courses given a prefix.
func (d *DefaultAutoCompleter) CoursesWithPrefix(prefix string) []string {
	return d.Courses.PrefixSearch(strings.ToUpper(prefix))
}
