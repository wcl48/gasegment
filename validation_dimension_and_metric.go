package gasegment

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/wcl48/gasegment/asset"
	analytics "google.golang.org/api/analytics/v3"
)

type DimensionOrMetricError string

func (e DimensionOrMetricError) Error() string { return string(e) }

var NoSuchDimensionOrMetric = DimensionOrMetricError("no such dimension or metric")

type ColumnAttributes struct {
	Id                      string
	ReplacedBy              string
	Type                    string
	DataType                string
	Group                   string
	Status                  string
	UIName                  string
	AppUIName               string
	Description             string
	Calculation             string
	MinTemplateIndex        int
	MaxTemplateIndex        int
	PremiumMinTemplateIndex int
	PremiumMaxTemplateIndex int
	AllowedInSegments       bool

	pattern *regexp.Regexp
}

func (ca *ColumnAttributes) Match(dm string) bool {
	if ca.pattern != nil {
		matches := ca.pattern.FindStringSubmatch(dm)
		if len(matches) < 2 {
			return false
		}
		digits := matches[1]
		if strings.HasPrefix(digits, "0") {
			return false
		}
		index, _ := strconv.Atoi(digits)
		return index >= ca.MinTemplateIndex && index <= ca.MaxTemplateIndex
	}
	return ca.Id == dm
}

var columns analytics.Columns
var dmDefMap = map[string]ColumnAttributes{}
var templateMatches []struct {
	Pattern          *regexp.Regexp
	ColumnAttributes *ColumnAttributes
}

func convertAttributes(id string, column *analytics.Column) ColumnAttributes {
	ca := ColumnAttributes{
		Id:                id,
		ReplacedBy:        column.Attributes["replacedBy"],
		Type:              column.Attributes["type"],
		DataType:          column.Attributes["dataType"],
		Group:             column.Attributes["group"],
		Status:            column.Attributes["status"],
		UIName:            column.Attributes["uiName"],
		AppUIName:         column.Attributes["appUiName"],
		Description:       column.Attributes["description"],
		Calculation:       column.Attributes["calculation"],
		AllowedInSegments: column.Attributes["allowedInSegments"] == "true",
	}

	ca.MinTemplateIndex, _ = strconv.Atoi(column.Attributes["minTemplateIndex"])
	ca.MaxTemplateIndex, _ = strconv.Atoi(column.Attributes["maxTemplateIndex"])
	ca.PremiumMinTemplateIndex, _ = strconv.Atoi(column.Attributes["premiumMinTemplateIndex"])
	ca.PremiumMaxTemplateIndex, _ = strconv.Atoi(column.Attributes["premiumMaxTemplateIndex"])

	return ca
}

func init() {
	b, err := asset.Asset("columns.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(b, &columns); err != nil {
		panic(err)
	}

	for _, column := range columns.Items {
		ca := convertAttributes(column.Id, column)
		if strings.Contains(ca.Id, "XX") {
			ca.pattern = regexp.MustCompile(strings.Replace(ca.Id, "XX", `(\d+)`, 1))
		}
		dmDefMap[column.Id] = ca
	}
}

func ColumnForDimensionOrMetric(dm string) (ColumnAttributes, error) {
	for _, ca := range dmDefMap {
		if ca.Match(dm) {
			return ca, nil
		}
	}

	return ColumnAttributes{}, NoSuchDimensionOrMetric
}
