# Google Analytics segment definition Parser and Serializer

![](https://codeship.com/projects/2d9e0530-881d-0132-48e0-16b8ca61b731/status?branch=master)

## Example

```go
package main

import (
	"fmt"

	"github.com/wacul/gasegment"
)

func main() {
	// parse
	segments, err := gasegment.Parse("users::condition::ga:pagePath==/abc")
	if err != nil {
		panic(err)
	}

	// modify
	segments[0].Scope = gasegment.SessionScope
	segments[0].Condition.AndExpression[0][0].Operator = gasegment.NotEqual
	segments[0].Condition.AndExpression[0][0].Value = "cde"

	// stringify
	fmt.Println(segments.DefString())
}
```

## Commandline

```
$ go get -v github.com/wacul/gasegment/cmd/gasegment
$ echo "sessions::condition::ga:medium==referral"  | gasegment
{
  "name": "-",
  "sessionSegment": {
    "segmentFilters": [
      {
        "simpleSegment": {
          "orFiltersForSegment": [
            {
              "segmentFilterClauses": [
                {
                  "dimensionFilter": {
                    "dimensionName": "ga:medium",
                    "expressions": [
                      "referral"
                    ],
                    "operator": "EXACT"
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
```
