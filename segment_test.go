package gasegment

import (
	"flag"
	"net/http"
	"reflect"
	"regexp"
	"testing"

	analytics "google.golang.org/api/analytics/v3"
)

func checkStringify(t *testing.T, set testset) {
	act := set.object.DefString()
	if act != set.definition {
		t.Errorf("failed to stringify\n\texpected: %s\n\tactual:   %s", set.definition, act)
	}
}

func checkParse(t *testing.T, set testset) {
	act, err := Parse(set.definition)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(act, set.object) {
		t.Errorf("failed to parse\nexpected: %v\nactual:   %v", set.object, act)
	}
}

func TestSegment(t *testing.T) {
	s := Segments{}
	if s.DefString() != "" {
		t.Errorf("empty Segments.DefString() must be \"\"")
	}

	for _, s := range set {
		checkStringify(t, s)
		checkParse(t, s)
	}

	for _, def := range checkDefs {
		s, err := Parse(def)
		if err != nil {
			t.Error(err)
			continue
		}
		act := s.DefString()
		if act != def {
			t.Errorf("check failed\n\texpected: %s\n\tactual:   %s", def, act)
		}
	}
}

func checkSplit(t *testing.T, s string, r *regexp.Regexp, expected []string) {
	act := splitByFirstRegexpGroup(s, r)
	if !reflect.DeepEqual(act, expected) {
		t.Errorf("splitByFirstRegexpGroup failed. actual: %v, expected: %v", act, expected)
	}
}

func TestSplitByFirstRegexpGroup(t *testing.T) {
	p := regexp.MustCompile(`aaa(bbb)ccc`)
	checkSplit(t, "abc", p, []string{"abc"})
	checkSplit(t, "hogehogeaaabbbcccfugafuga", p, []string{"hogehogeaaa", "bbb", "cccfugafuga"})
	checkSplit(t, "hogehogeaaabbbcccfugafugaaaabbbccc!!!", p, []string{"hogehogeaaa", "bbb", "cccfugafugaaaa", "bbb", "ccc!!!"})

	checkSplit(t, "a;->b;->>c", SequenceSeparatorRe, []string{"a", ";->", "b", ";->>", "c"})
}

func TestSortScope(t *testing.T) {
	input := `sessions::condition::ga:deviceCategory==desktop;users::condition::ga:pagePath!~^\Q/lk/\E,ga:pagePath!~^\Q/netacho/\E;sessions::condition::!ga:channelGrouping==(none);condition::ga:pagePath=~^\Q/inquiry/\E,ga:pagePath=~^\Q/inquiry/\E;condition::ga:deviceCategory==desktop;users::condition::ga:pagePath!~^\Q/lk/\E,ga:pagePath!~^\Q/netacho/\E;sessions::condition::!ga:channelGrouping==(none);condition::ga:goal4Completions>0`

	expected := `users::condition::ga:pagePath!~^\Q/lk/\E,ga:pagePath!~^\Q/netacho/\E;condition::ga:pagePath!~^\Q/lk/\E,ga:pagePath!~^\Q/netacho/\E;sessions::condition::ga:pagePath=~^\Q/inquiry/\E,ga:pagePath=~^\Q/inquiry/\E;condition::ga:deviceCategory==desktop;condition::!ga:channelGrouping==(none);condition::ga:goal4Completions>0;condition::ga:deviceCategory==desktop;condition::!ga:channelGrouping==(none)`

	ss, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}

	if ss.DefString() != expected {
		t.Error("bad sort")
	}
}

// Flags
var (
	testRequest = flag.Bool("runrequest", false, "Run test request to google analytics")
	clientID    = flag.String("clientid", "", "OAuth 2.0 Client ID.  If non-empty, overrides --clientid_file")
	secret      = flag.String("secret", "", "OAuth 2.0 Client Secret.  If non-empty, overrides --secret_file")
)

func TestValidSegment(t *testing.T) {
	flag.Parse()

	var client *http.Client
	var profileId string
	var svc *analytics.Service
	if *testRequest {

		if clientID == nil || secret == nil {
			panic("bad args")
			client = getTestClient(*clientID, *secret)

			svc, err := analytics.New(client)
			if err != nil {
				panic(err)
			}

			accountCall := svc.Management.AccountSummaries.List()
			accounts, err := accountCall.Do()
			if err != nil {
				panic(err)
			}

		LOOP:
			for _, account := range accounts.Items {
				for _, wp := range account.WebProperties {
					for _, pr := range wp.Profiles {
						profileId = pr.Id
						break LOOP
					}
				}
			}
		}
	}

	getSegmentCall := func(segments Segments) *analytics.DataGaGetCall {
		call := svc.Data.Ga.Get("ga:"+profileId, "yesterday", "today", "ga:sessions")
		call.SamplingLevel("FASTER")

		call.Segment(segments.DefString())
		return call
	}

	for _, valid := range valids {
		ss, err := Parse(valid)
		if err != nil {
			t.Fatal(err)
		}

		if client != nil {
			c := getSegmentCall(ss)
			if _, err := c.Do(); err != nil {
				t.Errorf("segment '%s' must suceed to call", valid)
			}
		}
	}

	for _, invalid := range invalids {
		ss, err := Parse(invalid)

		if err != nil {
			t.Fatal(err)
		}

		if client != nil {
			c := getSegmentCall(ss)

			_, err := c.Do()
			if err == nil {
				t.Errorf("segment '%s' must fail to call", invalid)
			}
		}
	}
}
