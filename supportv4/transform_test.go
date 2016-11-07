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

	// // https://developers.google.com/analytics/devguides/reporting/core/v4/migration
	// t.Run("1", func(t *testing.T) {
	// 	s := "users::condition::ga:userGender==Male;users::condition::ga:interestAffinityCategory==Games;sessions::condition::ga:region==Americas;sessions::condition::ga:language==en-u"
	// 	payload := stringToPayload(t, s)
	// 	fmt.Println(payload)
	// })
}
