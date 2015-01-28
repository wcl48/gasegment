package gasegment

import (
	"reflect"
	"regexp"
	"testing"
)

func checkStringify(t *testing.T, set testset) {
	act := set.object.DefString()
	if act != set.definition {
		t.Errorf("failed to stringify\n\texpected: %s\n\tactual:   %s", set.definition, act)
	}
}

func checkParse(t *testing.T, set testset) {
	act, err := Parse(set.definition)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(act, set.object) {
		t.Errorf("failed to parse\nexpected: %v\nactual:   %v", set.object, act)
	}
}

func TestSegment(t *testing.T) {
	s := Segments{}
	if s.DefString() != "" {
		t.Errorf("empty Segments.DefString() must be \"\"")
	}

	for _, s := range set {
		checkStringify(t, s)
		checkParse(t, s)
	}

	for _, def := range checkDefs {
		s, err := Parse(def)
		if err != nil {
			t.Error(err)
			continue
		}
		act := s.DefString()
		if act != def {
			t.Errorf("check failed\n\texpected: %s\n\tactual:   %s", def, act)
		}
	}
}

func checkSplit(t *testing.T, s string, r *regexp.Regexp, expected []string) {
	act := splitByFirstRegexpGroup(s, r)
	if !reflect.DeepEqual(act, expected) {
		t.Errorf("splitByFirstRegexpGroup failed. actual: %v, expected: %v", act, expected)
	}
}

func TestSplitByFirstRegexpGroup(t *testing.T) {
	p := regexp.MustCompile(`aaa(bbb)ccc`)
	checkSplit(t, "abc", p, []string{"abc"})
	checkSplit(t, "hogehogeaaabbbcccfugafuga", p, []string{"hogehogeaaa", "bbb", "cccfugafuga"})
	checkSplit(t, "hogehogeaaabbbcccfugafugaaaabbbccc!!!", p, []string{"hogehogeaaa", "bbb", "cccfugafugaaaa", "bbb", "ccc!!!"})

	checkSplit(t, "a;->b;->>c", SequenceSeparatorRe, []string{"a", ";->", "b", ";->>", "c"})
}
