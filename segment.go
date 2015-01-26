package gasegment

import (
	"fmt"
	"strings"
)

type DimensionOrMetric string

func (dm DimensionOrMetric) String() string {
	return string(dm)
}

type Operator string

func (op Operator) String() string {
	return string(op)
}

const (
	Equal                = Operator("==")
	NotEqual             = Operator("!=")
	LessThan             = Operator("<")
	LessThanEqual        = Operator("<=")
	GreaterThan          = Operator(">")
	GreaterThanEqual     = Operator(">=")
	Between              = Operator("<>")
	InList               = Operator("[]")
	ContainsSubstring    = Operator("=@")
	NotContainsSubstring = Operator("!@")
	Regexp               = Operator("=!")
	NotRegexp            = Operator("!~")
)

type ConditionScope string

func (cs ConditionScope) String() string {
	return string(cs)
}

const (
	UserScope    = ConditionScope("user::")
	SessionScope = ConditionScope("session::")
)

type ConditionType string

func (ct ConditionType) String() string {
	return string(ct)
}

const (
	ConditionSegment = ConditionType("condition::")
	SequenceSegment  = ConditionType("sequence::")
)

type MetricScope string

func (ms MetricScope) String() string {
	return string(ms)
}

const (
	Default    = MetricScope("")
	PerHit     = MetricScope("perHit::")
	PerSession = MetricScope("perSession::")
	PerUser    = MetricScope("perUser::")
)

type SequenceStepType string

func (st SequenceStepType) String() string {
	return string(st)
}

const (
	FirstStep        = SequenceStepType("")
	PrecedesSequence = SequenceStepType(";–>>")
	Immediately      = SequenceStepType(";–>")
)

type SegmentConditions []SegmentCondition

func NewSegmentConditions(cs ...SegmentCondition) SegmentConditions {
	return SegmentConditions(cs)
}

func (scs SegmentConditions) String() string {
	var currentScope ConditionScope
	buf := []string{}
	for _, sc := range scs {
		if currentScope != sc.Scope {
			buf = append(buf, sc.Scope.String()+sc.StringWithoutScope())
		} else {
			buf = append(buf, sc.StringWithoutScope())
		}
		currentScope = sc.Scope
	}
	return strings.Join(buf, ";")
}

type SegmentCondition struct {
	Scope        ConditionScope
	Type         ConditionType
	AndCondition AndCondition
	Sequence     *Sequence
}

func (sc *SegmentCondition) String() string {
	return sc.Scope.String() + sc.StringWithoutScope()
}

func (sc *SegmentCondition) StringWithoutScope() string {
	switch sc.Type {
	case ConditionSegment:
		return sc.Type.String() + sc.AndCondition.String()
	case SequenceSegment:
		return sc.Type.String() + sc.Sequence.String()
	default:
		return ""
	}
}

type AndCondition struct {
	Exclude      bool
	OrConditions []OrCondition
}

func (a AndCondition) String() string {
	buf := []string{}
	if a.Exclude {
		buf = append(buf, "!")
	}
	for _, or := range a.OrConditions {
		buf = append(buf, or.String())
	}
	return strings.Join(buf, ";")
}

type OrCondition struct {
	Conditions []Condition
}

func (o OrCondition) String() string {
	buf := make([]string, len(o.Conditions))
	for i, condition := range o.Conditions {
		buf[i] = condition.String()
	}
	return strings.Join(buf, ",")
}

type Condition struct {
	MetricScope MetricScope
	Target      DimensionOrMetric
	Operator    Operator
	Value       string
}

func (c Condition) EscapedValue() string {
	return strings.Replace(strings.Replace(c.Value, ";", "\\;", -1), ",", "\\,", -1)
}

func (c Condition) String() string {
	return stringerJoin("", c.MetricScope, c.Target, c.Operator) + c.EscapedValue()
}

type Sequence struct {
	Not                      bool
	FirstHitMatchesFirstStep bool
	SequensSteps             SequenceSteps
}

func (s Sequence) String() string {
	var buf [3]string
	if s.Not {
		buf[0] = "!"
	}
	if s.FirstHitMatchesFirstStep {
		buf[1] = "^"
	}
	buf[2] = s.SequensSteps.String()
	return strings.Join(buf[:], "")
}

type SequenceStep struct {
	Type         SequenceStepType
	AndCondition AndCondition
}

type SequenceSteps []SequenceStep

func (ss SequenceSteps) String() string {
	buf := make([]string, len((ss)))
	for i, s := range ss {
		buf[i] = s.String()
	}
	return strings.Join(buf, "")
}

func (ss SequenceStep) String() string {
	return stringerJoin("", ss.Type, ss.AndCondition)
}

func stringerJoin(sep string, stringers ...fmt.Stringer) string {
	buf := make([]string, len(stringers))
	for i, stringer := range stringers {
		buf[i] = stringer.String()
	}
	return strings.Join(buf, sep)
}
