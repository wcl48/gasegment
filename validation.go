package gasegment

import "unicode/utf8"

func (ss Segments) IsValid() bool {
	for _, s := range ss {
		if !s.IsValid() {
			return false
		}
	}
	return true
}

func (s Segment) IsValid() bool {
	switch s.Type {
	case ConditionSegment:
		return s.Condition.AndExpression.IsValid()
	case SequenceSegment:

	}
	return true
}

func (ae AndExpression) IsValid() bool {
	for _, oe := range ae {
		if !oe.IsValid() {
			return false
		}
	}

	return true
}

func (oe OrExpression) IsValid() bool {
	for _, e := range oe {
		if !e.IsValid() {
			return false
		}
	}

	return true
}

func (e Expression) IsValid() bool {
	switch e.Operator {
	case Regexp, NotRegexp:
		if utf8.RuneCountInString(e.Value) > 128 {
			return false
		}
	default:
		if utf8.RuneCountInString(e.Value) > 1024 {
			return false
		}
	}

	return true
}
