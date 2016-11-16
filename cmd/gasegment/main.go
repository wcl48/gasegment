package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/wcl48/gasegment"
	"github.com/wcl48/gasegment/supportv4"
	"google.golang.org/api/analyticsreporting/v4"
)

func parse(reader io.Reader) (*analyticsreporting.DynamicSegment, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	defstring := strings.Trim(string(b), "\n")
	segments, err := gasegment.Parse(defstring)
	if err != nil {
		return nil, err
	}
	ds, err := supportv4.TransformSegments(&segments)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

func dump(ds *analyticsreporting.DynamicSegment) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(ds)
}

func main() {
	if len(os.Args) < 2 {
		ds, err := parse(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		dump(ds)
	} else {
		for _, fname := range os.Args[1:] {
			f, err := os.Open(fname)
			defer f.Close()
			if err != nil {
				log.Fatal(err)
			}
			ds, err := parse(f)
			if err != nil {
				log.Fatal(err)
			}
			dump(ds)

		}
	}
}
