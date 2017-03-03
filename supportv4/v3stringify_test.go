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
		"sessions::condition::ga:sessionCount>2;ga:sessionCount<>foo\\_bar_end",
		"users::sequence::ga:deviceCategory==desktop;->ga:deviceCategory==tablet",
		"sessions::condition::ga:medium=~^(cpc|ppc|cpa|cpm|cpv|cpp|xx\\|xx)$",
		"users::condition::!ga:pagePath=~^\\Q/recruit/\\E;sessions::condition::ga:deviceCategory=@desktop",
		"users::condition::perSession::ga:goal3Completions!=0;condition::!ga:pagePath=~^\\Q/recruit/\\E;condition::ga:pagePath=~^\\Q/kpdf_app\\E;sessions::condition::ga:deviceCategory=@desktop",
		"users::sequence::!ga:flashVersion=@bar2;sessions::condition::ga:flashVersion=@test\\,ga:flashVersion=@test3\\\\;ga:flashVersion=@and\\;\\,;condition::!ga:sessionDurationBucket==123;sequence::^ga:operatingSystem=@Windows",
		"sessions::condition::ga:deviceCategory=@mobile;condition::ga:landingPagePath==sp.exampleonline.co.jp/exampleol/www.exampleonline.co.jp/mb/BSfMbCategoryTop.jsp\\;sjid=CB706677EAC2E584171A5582CC8275C6.c?bg=set&guid=on",
		"sessions::condition::ga:deviceCategory=@desktop;condition::ga:pagePath=@embed,ga:pagePath==/files/embed/cartonbox.html,ga:pagePath=@/files/cp/kaitori,ga:pagePath=@/cd/files/kaitori1307,ga:pagePath==/files/selltop.html;ga:pagePath=@sell,ga:pagePath==/files/embed/cartonbox.html,ga:pagePath=@/files/cp/kaitori,ga:pagePath=@/cd/files/kaitori1307,ga:pagePath==/files/selltop.html",
		"users::sequence::!^ga:pagePath==/aiueo;->ga:pagePath==/aiueo2;->>ga:pagePath==/aiueo3",
		"users::sequence::!^ga:pagePath==/aiueo;->>ga:pagePath==/aiueo2;->ga:pagePath==/aiueo3",
		"users::sequence::^ga:sessionCount==1;dateOfSession<>2014-05-20_2014-05-30;->>ga:sessionDurationBucket>600",
		"sessions::condition::ga:pagePath"
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

func TestReversible2(t *testing.T) {
	for i, defstring := range gasegment.TestCheckDefs {
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
