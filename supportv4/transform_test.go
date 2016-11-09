package supportv4

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/wcl48/gasegment"
)

func TestTransformNotError(t *testing.T) {
	// anyway, checking that all test set data are transformed for v4.
	for _, s := range gasegment.TestCheckDefs {
		t.Run(s, func(t *testing.T) {
			segments, err := gasegment.Parse(s)
			if err != nil {
				t.Fatalf("error raised: %s\n", err)
			}
			ds, err := TransformSegments(&segments)
			if err != nil {
				t.Fatalf("error raised: %s\n", err)
			}
			_, err = json.MarshalIndent(ds, "", "  ")
			if err != nil {
				t.Fatalf("error raised: %s\n", err)
			}
		})
	}
}

func assertJSONEqual(t *testing.T, expectedJSON string, actualJSON string) {
	var ob interface{}
	var ob2 interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &ob); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(actualJSON), &ob2); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ob, ob2) {
		t.Errorf("expected=%q, but actual=%q", expectedJSON, actualJSON)
	}
}

func TestTransform(t *testing.T) {
	stringToPayload := func(t *testing.T, s string, name string) string {
		segments, err := gasegment.Parse(s)
		if err != nil {
			t.Fatalf("error raised: %s\n", err)
		}
		ds, err := TransformSegments(&segments)
		if err != nil {
			t.Fatalf("error raised: %s\n", err)
		}
		ds.Name = name
		b, err := json.MarshalIndent(ds, "", "  ")
		if err != nil {
			t.Fatalf("error raised: %s\n", err)
		}
		return string(b)
	}

	// from: https://developers.google.com/analytics/devguides/reporting/core/v4/migration
	t.Run("dynamic segment 1(migration)", func(t *testing.T) {
		s := "sessions::condition::ga:medium==referral"
		expectedJSON := `
{
  "name": "segment_name",
  "sessionSegment": {
    "segmentFilters": [{
      "simpleSegment": {
        "orFiltersForSegment": [{
          "segmentFilterClauses": [{
            "dimensionFilter": {
              "dimensionName": "ga:medium",
              "operator": "EXACT",
              "expressions": [ "referral" ]
            }
          }]
        }]
      }
    }]
  }
}
`
		transformedJSON := stringToPayload(t, s, "segment_name")
		assertJSONEqual(t, expectedJSON, transformedJSON)
	})
	t.Run("dynamic segment 2(migration)", func(t *testing.T) {
		// original := "users::condition::ga:revenue>10;sequence::ga:deviceCategory==desktop->>ga:deviceCategory==mobile"
		// `ga:revenue` is not found. so `ga:sessions` is used.
		s := "users::condition::ga:sessions>10;sequence::ga:deviceCategory==desktop;->>ga:deviceCategory==mobile"
		expectedJSON := `
{
  "name": "segment_name",
  "userSegment": {
    "segmentFilters": [
      {
        "simpleSegment": {
          "orFiltersForSegment": [{
            "segmentFilterClauses": [{
              "metricFilter": {
                "metricName": "ga:sessions",
                "operator": "GREATER_THAN",
                "comparisonValue": "10"
              }
            }]
          }]
        }
      },
      {
        "sequenceSegment": {
          "segmentSequenceSteps": [{
            "orFiltersForSegment": [{
              "segmentFilterClauses": [{
                "dimensionFilter": {
                  "dimensionName": "ga:deviceCategory",
                  "operator": "EXACT",
                  "expressions": ["desktop"]
                }
              }]
            }]
          },{
            "matchType": "PRECEDES",
            "orFiltersForSegment": [{
              "segmentFilterClauses": [{
                "dimensionFilter": {
                  "dimensionName": "ga:deviceCategory",
                  "operator": "EXACT",
                  "expressions": ["mobile"]
                }
              }]
            }]
          }]
        }
      }
    ]
  }
}
`
		transformedJSON := stringToPayload(t, s, "segment_name")
		assertJSONEqual(t, expectedJSON, transformedJSON)
	})

	t.Run("appendix 1", func(t *testing.T) {
		s := "sessions::condition::!ga:landingPagePath=~^\\Qexample.com/blog/xxx/\\E;condition::!ga:landingPagePath=~^\\Qexample.com/yyy/\\E"
		expectedJSON := `
{
  "name": "segment_name",
  "sessionSegment": {
    "segmentFilters": [
      {
        "simpleSegment": {
          "orFiltersForSegment": [
            {
              "segmentFilterClauses": [
                {
                  "dimensionFilter": {
                    "expressions": [
                      "^\\Qexample.com/blog/xxx/\\E"
                    ],
                    "dimensionName": "ga:landingPagePath",
                    "operator": "REGEXP"
                  }
                }
              ]
            }
          ]
        },
        "not": true
      },
      {
        "simpleSegment": {
          "orFiltersForSegment": [
            {
              "segmentFilterClauses": [
                {
                  "dimensionFilter": {
                    "expressions": [
                      "^\\Qexample.com/yyy/\\E"
                    ],
                    "dimensionName": "ga:landingPagePath",
                    "operator": "REGEXP"
                  }
                }
              ]
            }
          ]
        },
        "not": true
      }
    ]
  }
}
`
		transformedJSON := stringToPayload(t, s, "segment_name")
		assertJSONEqual(t, expectedJSON, transformedJSON)
	})

	t.Run("appendix 2 (combined with regexp)", func(t *testing.T) {
		s := "sessions::condition::ga:landingPagePath=~^example.com/blog/(xxx|yyy)/"
		expectedJSON := `
{
  "name": "segment_name",
  "sessionSegment": {
    "segmentFilters": [
      {
        "simpleSegment": {
          "orFiltersForSegment": [
            {
              "segmentFilterClauses": [
                {
                  "dimensionFilter": {
                    "expressions": [
                      "^example.com/blog/(xxx|yyy)/"
                    ],
                    "dimensionName": "ga:landingPagePath",
                    "operator": "REGEXP"
                  }
                }
              ]
            }
          ]
        }
      }
    ]
  }
}
`
		transformedJSON := stringToPayload(t, s, "segment_name")
		assertJSONEqual(t, expectedJSON, transformedJSON)
	})

	t.Run("appendix 2 (combined with or operator)", func(t *testing.T) {
		s := "sessions::condition::ga:landingPagePath=~^\\Qexample.com/blog/xxx\\E,ga:landingPagePath=~^\\Qexample.com/blog/yyy\\E"
		expectedJSON := `
{
  "name": "segment_name",
  "sessionSegment": {
    "segmentFilters": [
      {
        "simpleSegment": {
          "orFiltersForSegment": [
            {
              "segmentFilterClauses": [
                {
                  "dimensionFilter": {
                    "operator": "REGEXP",
                    "expressions": [
                      "^\\Qexample.com/blog/xxx\\E"
                    ],
                    "dimensionName": "ga:landingPagePath"
                  }
                },
                {
                  "dimensionFilter": {
                    "operator": "REGEXP",
                    "expressions": [
                      "^\\Qexample.com/blog/yyy\\E"
                    ],
                    "dimensionName": "ga:landingPagePath"
                  }
                }
              ]
            }
          ]
        }
      }
    ]
  }
}
`
		transformedJSON := stringToPayload(t, s, "segment_name")
		assertJSONEqual(t, expectedJSON, transformedJSON)

	})

	t.Run("appendix 3 greater than and between", func(t *testing.T) {
		s := "sessions::condition::ga:sessionCount>2;ga:sessionCount<>2_3;ga:hits>10;ga:hits<>10_100"
		expectedJSON := `
{
  "name": "segment_name",
  "sessionSegment": {
    "segmentFilters": [
      {
        "simpleSegment": {
          "orFiltersForSegment": [
            {
              "segmentFilterClauses": [
                {
                  "dimensionFilter": {
                    "dimensionName": "ga:sessionCount",
                    "operator": "NUMERIC_GREATER_THAN",
                    "expressions": [
                      "2"
                    ]
                  }
                }
              ]
            },
            {
              "segmentFilterClauses": [
                {
                  "dimensionFilter": {
                    "maxComparisonValue": "3",
                    "dimensionName": "ga:sessionCount",
                    "operator": "NUMERIC_BETWEEN",
                    "minComparisonValue": "2"
                  }
                }
              ]
            },
            {
              "segmentFilterClauses": [
                {
                  "metricFilter": {
                    "comparisonValue": "10",
                    "operator": "GREATER_THAN",
                    "metricName": "ga:hits"
                  }
                }
              ]
            },
            {
              "segmentFilterClauses": [
                {
                  "metricFilter": {
                    "comparisonValue": "10",
                    "maxComparisonValue": "100",
                    "operator": "BETWEEN",
                    "metricName": "ga:hits"
                  }
                }
              ]
            }
          ]
        }
      }
    ]
  }
}
`
		transformedJSON := stringToPayload(t, s, "segment_name")
		assertJSONEqual(t, expectedJSON, transformedJSON)

	})

}

func TestBrokenBetween(t *testing.T) {
	t.Run("on metric", func(t *testing.T) {
		invalid := "sessions::condition::ga:hits<>10" // not ga:hits<>10_20
		segments, err := gasegment.Parse(invalid)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := TransformSegments(&segments); err == nil {
			t.Error("must be error")
		}
	})
	t.Run("on dimension", func(t *testing.T) {
		invalid := "sessions::condition::ga:sessionCount<>10" // not ga:hits<>10_20
		segments, err := gasegment.Parse(invalid)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := TransformSegments(&segments); err == nil {
			t.Error("must be error")
		}
	})
}
