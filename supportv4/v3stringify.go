package supportv4

import (
	"errors"
	"fmt"
	"strings"

	gapi "google.golang.org/api/analyticsreporting/v4"
)

// V4[AST] -> V3[string]

// V3StringifySegmentFilterClause :
func V3StringifySegmentFilterClause(node *gapi.SegmentFilterClause) (string, error) {
	if node.DimensionFilter != nil {
		return V3StringifySegmentDimensionFilter(node.DimensionFilter, node.Not)
	}
	if node.MetricFilter != nil {
		return V3StringifySegmentMetricFilter(node.MetricFilter, node.Not)
	}
	return "", errors.New("must be wither a metric or a dimension filter")
}

// V3StringifySegmentDimensionFilter :
func V3StringifySegmentDimensionFilter(node *gapi.SegmentDimensionFilter, not bool) (string, error) {
	if node == nil {
		return "", nil
	}

	// "OPERATOR_UNSPECIFIED" - If the match type is unspecified, it is treated as a REGEXP.
	op := node.Operator
	if op == "OPERATOR_UNSPECIFIED" || op == "" {
		op = OperatorRegexp
	}
	if len(node.Expressions) == 0 {
		return "", errors.New("invalid expression. at least length >= 1")
	}
	// see also: ./detect.go DetectOperatorOnDimension

	// TODO: node.CaseSensitive
	switch op {
	case OperatorRegexp:
		if not {
			return fmt.Sprintf("%s!~%s", node.DimensionName, node.Expressions[0]), nil
		}
		return fmt.Sprintf("%s=~%s", node.DimensionName, node.Expressions[0]), nil
	case OperatorBeginsWith:
		// TODO: need regex.QuoteMeta() ?
		if not {
			return fmt.Sprintf("%s!~%s", node.DimensionName, "^"+node.Expressions[0]), nil
		}
		return fmt.Sprintf("%s=~%s", node.DimensionName, "^"+node.Expressions[0]), nil
	case OperatorEndsWith:
		if not {
			return fmt.Sprintf("%s!~%s", node.DimensionName, node.Expressions[0]+"$"), nil
		}
		return fmt.Sprintf("%s=~%s", node.DimensionName, node.Expressions[0]+"$"), nil
	case OperatorPartial:
		if not {
			return fmt.Sprintf("%s!@%s", node.DimensionName, node.Expressions[0]), nil
		}
		return fmt.Sprintf("%s=@%s", node.DimensionName, node.Expressions[0]), nil
	case OperatorExact, OperatorNumericEquals:
		if not {
			return fmt.Sprintf("%s!=%s", node.DimensionName, node.Expressions[0]), nil
		}
		return fmt.Sprintf("%s==%s", node.DimensionName, node.Expressions[0]), nil
	case OperatorInList:
		// TODO: limitation of number of expressions <= 10
		// TODO: need escape?
		if not {
			return "", fmt.Errorf("not support %q with Not=true", OperatorInList)
		}
		return fmt.Sprintf("%s[]%s", node.DimensionName, strings.Join(node.Expressions, "|")), nil
	case OperatorNumericLessThan:
		if not {
			return fmt.Sprintf("%s>=%s", node.DimensionName, node.Expressions[0]), nil
		}
		return fmt.Sprintf("%s<%s", node.DimensionName, node.Expressions[0]), nil
	case OperatorNumericGreaterThan:
		if not {
			return fmt.Sprintf("%s<=%s", node.DimensionName, node.Expressions[0]), nil
		}
		return fmt.Sprintf("%s>%s", node.DimensionName, node.Expressions[0]), nil
	case OperatorNumericBetween:
		if not {
			return "", fmt.Errorf("not support %q with Not=true", OperatorNumericBetween)
		}
		return fmt.Sprintf("%s<>%s_%s", node.DimensionName, node.Expressions[0], node.Expressions[1]), nil
	default:
		return "", fmt.Errorf("unsupported dimension operator: %s", op)
	}
}

// V3StringifySegmentMetricFilter :
func V3StringifySegmentMetricFilter(node *gapi.SegmentMetricFilter, not bool) (string, error) {
	if node == nil {
		return "", nil
	}

	// TODO: node.Scope
	// "OPERATOR_UNSPECIFIED" - If the match type is unspecified, it is treated as a REGEXP.
	op := node.Operator
	if op == "OPERATOR_UNSPECIFIED" {
		op = OperatorEqual
	}
	// see also: ./detect.go DetectOperatorOnMetric

	scopePrefix := ""
	switch node.Scope {
	case "PRODUCT":
		scopePrefix = "perProduct::" // Product scope.
	case "HIT":
		scopePrefix = "perHit::" // Hit scope.
	case "USER":
		scopePrefix = "perUser::" // User scope.
	case "SESSION":
		scopePrefix = "perSession::" // Session scope.
	}
	switch op {
	case OperatorEqual:
		if not {
			return fmt.Sprintf("%s%s!=%s", scopePrefix, node.MetricName, node.ComparisonValue), nil
		}
		return fmt.Sprintf("%s%s==%s", scopePrefix, node.MetricName, node.ComparisonValue), nil
	case OperatorLessThan:
		if not {
			return fmt.Sprintf("%s%s>=%s", scopePrefix, node.MetricName, node.ComparisonValue), nil
		}
		return fmt.Sprintf("%s%s<%s", scopePrefix, node.MetricName, node.ComparisonValue), nil
	case OperatorGreaterThan:
		if not {
			return fmt.Sprintf("%s%s<=%s", scopePrefix, node.MetricName, node.ComparisonValue), nil
		}
		return fmt.Sprintf("%s%s>%s", scopePrefix, node.MetricName, node.ComparisonValue), nil
	case OperatorBetween:
		if not {
			return "", fmt.Errorf("not support %q with Not=true", OperatorBetween)
		}
		return fmt.Sprintf("%s%s<>%s_%s", scopePrefix, node.MetricName, node.ComparisonValue, node.MaxComparisonValue), nil
	default:
		return "", fmt.Errorf("unsupported metric operator: %s", op)
	}
}
