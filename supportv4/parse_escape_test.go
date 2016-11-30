package supportv4

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseInListValue(t *testing.T) {
	candidates := []struct {
		value    string
		expected []string
	}{
		{value: "foo bar", expected: []string{"foo bar"}},
		{value: "foo|bar", expected: []string{"foo", "bar"}},
		{value: "foo\\|bar", expected: []string{"foo|bar"}},
		{value: "foo\\\\|bar", expected: []string{"foo\\", "bar"}},
		{value: "foo\a\\|bar", expected: []string{"foo\a|bar"}},
	}

	for _, c := range candidates {
		t.Run(fmt.Sprintf("parseInList %s", c.value), func(t *testing.T) {
			parsed := ParseStringWithEscape(c.value, '|', '\\')
			if !reflect.DeepEqual(c.expected, parsed) {
				t.Errorf("expected %v, but %v", c.expected, parsed)
			}
		})
	}
}
