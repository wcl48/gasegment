package gasegment

import "testing"

func TestValidateDimensionOrMetric(t *testing.T) {
	table := []struct {
		dm    string
		id    string
		valid bool
	}{
		// valids
		{"ga:sessions", "ga:sessions", true},
		{"ga:goal12Starts", "ga:goalXXStarts", true},
		{"ga:goal1Starts", "ga:goalXXStarts", true},

		// invalids
		{"ga:session2", "", false},
		{"ga:goal0Starts", "", false},
		{"ga:goal02Starts", "", false},
		{"ga:goal21Starts", "", false},
	}

	for _, pattern := range table {
		ca, err := ColumnForDimensionOrMetric(pattern.dm)
		if pattern.valid {
			if err != nil {
				t.Errorf("unexpected error for %s : %s", pattern.dm, err.Error())
			}
			if pattern.id != ca.Id {
				t.Errorf("unexpected dimension for %s : expected %s, actual %s", pattern.dm, pattern.id, ca.Id)
			}
		} else {
			if err != NoSuchDimensionOrMetric {
				t.Errorf("unexpected error for %s : %s", pattern.dm, err.Error())
			}
		}
	}
}
