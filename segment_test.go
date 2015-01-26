package gasegment

import (
	"testing"
)

type testset struct {
	original string
	expected string
	object   SegmentConditions
}

var set = []testset{{
	original: "user::condition::ga:pagePath==/aiueo",
	expected: "user::condition::ga:pagePath==/aiueo",
	object: NewSegmentConditions(SegmentCondition{
		Scope: UserScope,
		Type:  ConditionSegment,
		AndCondition: AndCondition{
			OrConditions: []OrCondition{{
				Conditions: []Condition{{
					Target:   DimensionOrMetric("ga:pagePath"),
					Operator: Equal,
					Value:    "/aiueo",
				}},
			}},
		},
	}),
}, {
	original: "user::condition::ga:pagePath==/aiueo",
	expected: "user::condition::ga:pagePath==/aiueo;ga:pagePath==/bcdef\\;\\,",
	object: NewSegmentConditions(SegmentCondition{
		Scope: UserScope,
		Type:  ConditionSegment,
		AndCondition: AndCondition{
			OrConditions: []OrCondition{{
				Conditions: []Condition{{
					Target:   DimensionOrMetric("ga:pagePath"),
					Operator: Equal,
					Value:    "/aiueo",
				}},
			}, {
				Conditions: []Condition{{
					Target:   DimensionOrMetric("ga:pagePath"),
					Operator: Equal,
					Value:    "/bcdef;,",
				}},
			}},
		},
	}),
}, {
	original: "user::condition::ga:pagePath==/aiueo",
	expected: "user::condition::ga:pagePath==/aiueo;condition::ga:landingPagePath==/123",
	object: NewSegmentConditions(SegmentCondition{
		Scope: UserScope,
		Type:  ConditionSegment,
		AndCondition: AndCondition{
			OrConditions: []OrCondition{{
				Conditions: []Condition{{
					Target:   DimensionOrMetric("ga:pagePath"),
					Operator: Equal,
					Value:    "/aiueo",
				}},
			}},
		},
	}, SegmentCondition{
		Scope: UserScope,
		Type:  ConditionSegment,
		AndCondition: AndCondition{
			OrConditions: []OrCondition{{
				Conditions: []Condition{{
					Target:   DimensionOrMetric("ga:landingPagePath"),
					Operator: Equal,
					Value:    "/123",
				}},
			}},
		},
	}),
}}

func checkStringify(t *testing.T, set testset) {
	act := set.object.String()
	if act != set.expected {
		t.Errorf("stringify\n\texpected: %s\n\tactual:   %s", set.expected, act)
	}
}

func TestSegment(t *testing.T) {
	// stringify
	for _, s := range set {
		checkStringify(t, s)
	}
}
