package supportv4

import (
	"fmt"
	"testing"

	"github.com/wcl48/gasegment"
)

func fullTransform(defstring string) (string, error) {
	// v3string -> gasegment -> v4ast -> v3string
	segments, err := gasegment.Parse(defstring)
	if err != nil {
		return "", err
	}
	ds, err := TransformSegments(&segments)
	if err != nil {
		return "", err
	}
	return V3StringifyDynamicSegment(ds)
}

func TestReversible(t *testing.T) {
	candidates := []string{
		"sessions::condition::ga:medium==referral",
		"users::condition::ga:sessions>10;sequence::ga:deviceCategory==desktop;->>ga:deviceCategory==mobile",
		"sessions::condition::!ga:landingPagePath=~^\\Qexample.com/blog/xxx/\\E;condition::!ga:landingPagePath=~^\\Qexample.com/yyy/\\E",
		"sessions::condition::ga:landingPagePath=~^example.com/blog/(xxx|yyy)/",
		"sessions::condition::ga:landingPagePath=~^\\Qexample.com/blog/xxx\\E,ga:landingPagePath=~^\\Qexample.com/blog/yyy\\E",
		"sessions::condition::ga:sessionCount>2;ga:sessionCount<>2_3;ga:hits>10;ga:hits<>10_100",
		"users::sequence::ga:deviceCategory==desktop;->ga:deviceCategory==tablet",
	}

	for i, defstring := range candidates {
		defstring := defstring
		t.Run(fmt.Sprintf("reversible%d", i), func(t *testing.T) {
			result, err := fullTransform(defstring)
			if err != nil {
				t.Fatalf("no error required in %q. but %v is occured", defstring, err)
			}
			if result != defstring {
				t.Errorf("\nexpected:\n\t%q\ngot:\n\t%q", defstring, result)
			}
		})
	}
}
