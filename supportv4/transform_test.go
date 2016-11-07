package supportv4

import (
	"encoding/json"
	"testing"

	"github.com/wcl48/gasegment"
)

func TestTransform(t *testing.T) {
	// anyway, checking that all test set data are transformed to v4 version.
	for _, s := range gasegment.TestCheckDefs {
		t.Run(s, func(t *testing.T) {
			segments, err := gasegment.Parse(s)
			if err != nil {
				t.Fatal(err)
			}
			ds, err := TransformSegments(&segments)
			if err != nil {
				t.Fatal(err)
			}
			_, err = json.MarshalIndent(ds, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			// t.Log(string(b))
		})
	}
}
