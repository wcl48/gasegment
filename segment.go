package gasegment

import (
	"sort"
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
	Regexp               = Operator("=~")
	NotRegexp            = Operator("!~")
)

type SegmentScope string

func (cs SegmentScope) String() string {
	return string(cs)
}

const (
	UserScope    = SegmentScope("users::")
	SessionScope = SegmentScope("sessions::")
)

type SegmentType string

func (ct SegmentType) String() string {
	return string(ct)
}

const (
	ConditionSegment = SegmentType("condition::")
	SequenceSegment  = SegmentType("sequence::")
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
	FirstStep           = SequenceStepType("")
	Precedes            = SequenceStepType(";->>")
	ImmediatelyPrecedes = SequenceStepType(";->")
)

type Segments []Segment

var scopeSortMap = map[SegmentScope]int{
	UserScope:    0,
	SessionScope: 1,
}

type sortByScope []Segment

func (s sortByScope) Len() int {
	return len(s)
}

func (s sortByScope) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortByScope) Less(i, j int) bool {
	return scopeSortMap[s[i].Scope] < scopeSortMap[s[j].Scope]
}

func NewSegments(cs ...Segment) Segments {
	return Segments(cs)
}

func (scs Segments) DefString() string {
	workSegments := make([]Segment, len(scs))
	copy(workSegments, scs)

	sort.Sort(sortByScope(workSegments))

	var currentScope SegmentScope
	buf := []string{}
	for _, sc := range workSegments {
		scDef := sc.DefStringWithoutScope()
		if scDef == "" {
			continue
		}
		if currentScope != sc.Scope {
			buf = append(buf, sc.Scope.String()+scDef)
		} else {
			buf = append(buf, scDef)
		}
		currentScope = sc.Scope
	}
	return strings.Join(buf, ";")
}

func (scs *Segments) AddSegment(ss ...Segment) {
	*scs = append(*scs, ss...)
}

func (scs *Segments) AddSegments(sgs ...Segments) {
	for _, sg := range sgs {
		scs.AddSegment(sg...)
	}
}

type Segment struct {
	Scope     SegmentScope
	Type      SegmentType
	Condition Condition
	Sequence  Sequence
}

func (sc *Segment) DefString() string {
	return sc.Scope.String() + sc.DefStringWithoutScope()
}

func (sc *Segment) DefStringWithoutScope() string {
	switch sc.Type {
	case ConditionSegment:
		return sc.Type.String() + sc.Condition.DefString()
	case SequenceSegment:
		return sc.Type.String() + sc.Sequence.DefString()
	default:
		return ""
	}
}

type Condition struct {
	Exclude       bool
	AndExpression AndExpression
}

func (c Condition) DefString() string {
	buf := []string{}
	if c.Exclude {
		buf = append(buf, "!")
	}
	buf = append(buf, c.AndExpression.DefString())
	return strings.Join(buf, "")
}

type AndExpression []OrExpression

func (a AndExpression) DefString() string {
	buf := []string{}
	for _, or := range a {
		buf = append(buf, or.DefString())
	}
	return strings.Join(buf, ";")
}

func NewAndExpression(inner ...OrExpression) AndExpression {
	return AndExpression(inner)
}

func NewSingleAndExpression(es ...Expression) AndExpression {
	return NewAndExpression(NewOrExpression(es...))
}

type OrExpression []Expression

func (o OrExpression) DefString() string {
	buf := make([]string, len(o))
	for i, condition := range o {
		buf[i] = condition.DefString()
	}
	return strings.Join(buf, ",")
}

func NewOrExpression(inner ...Expression) OrExpression {
	return OrExpression(inner)
}

type Expression struct {
	MetricScope MetricScope
	Target      DimensionOrMetric
	Operator    Operator
	Value       string
}

func (c Expression) EscapedValue() string {
	return EscapeExpressionValue(c.Value)
}

func EscapeExpressionValue(s string) string {
	return strings.Replace(strings.Replace(s, ";", `\;`, -1), ",", `\,`, -1)
}
func UnEscapeExpressionValue(s string) string {
	return strings.Replace(strings.Replace(s, `\;`, ";", -1), `\,`, ",", -1)
}

func (c Expression) DefString() string {
	return strings.Join([]string{c.MetricScope.String(), c.Target.String(), c.Operator.String()}, "") + c.EscapedValue()
}

type Sequence struct {
	Not                      bool
	FirstHitMatchesFirstStep bool
	SequenceSteps            SequenceSteps
}

func (s Sequence) DefString() string {
	var buf [3]string
	if s.Not {
		buf[0] = "!"
	}
	if s.FirstHitMatchesFirstStep {
		buf[1] = "^"
	}
	buf[2] = s.SequenceSteps.DefString()
	return strings.Join(buf[:], "")
}

type SequenceStep struct {
	Type          SequenceStepType
	AndExpression AndExpression
}

type SequenceSteps []SequenceStep

func NewSequenceSteps(inner ...SequenceStep) SequenceSteps {
	return SequenceSteps(inner)
}

func (ss SequenceSteps) DefString() string {
	buf := make([]string, len((ss)))
	for i, s := range ss {
		buf[i] = s.DefString()
	}
	return strings.Join(buf, "")
}

func (ss SequenceStep) DefString() string {
	return strings.Join([]string{ss.Type.String(), ss.AndExpression.DefString()}, "")
}
