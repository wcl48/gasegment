package gasegment

import (
	analyticsreporting "google.golang.org/api/analyticsreporting/v4"
)

// transform gasegment.Segment -> analyticsreporting.DynamicSegment
// go get -v google.golang.org/api/analyticsreporting/v4
// (https://github.com/google/google-api-go-client/blob/master/analyticsreporting/v4/analyticsreporting-gen.go)

/*
# gaesegment

Segment
  Scope: SegmentScope{UserScope, SessionScope}
  Type: SegmentType{ConditionSegment, SequenceSegment}
  Condition:
    Exclude: bool{true/false}
    AndExpression as OrExpression[]:
      Expression[]
        * MetricScope: MetricScope{Default,PerHit,PerSession,PerUser}
        Target: DimensionOrMetric{}
        Operator: Operator{Equal,NotEqual,LessThan,LessThanEqual,GreaterThan,GreaterThanEqual,Between,InList,ContainsSubstring,NotContainsSubstring,Regexp,NotRegexp}
        Value: string{}
  Sequence:
    Not: bool{true/false}
    FirstHitMatchesFirstStep: bool{true/false}
    SequenceSteps[]:
      Type: SequenceStepType{FirstStep,Precedes,ImmediatelyPrecedes}
      AndExpression as OrExpression[]:
        Expression[]
          MetricScope: MetricScope{Default,PerHit,PerSession,PerUser}
          Target: DimensionOrMetric{}
          Operator: Operator{Equal,NotEqual,LessThan,LessThanEqual,GreaterThan,GreaterThanEqual,Between,InList,ContainsSubstring,NotContainsSubstring,Regexp,NotRegexp}
          Value: string{}

# analyticsreporting

DynamicSegment
  Name: string{}
  SessionSegment:
    SegmentFilters[]:
      Not: bool{true/false}
      SequenceSegment:
        FirstStepShouldMatchFirstHit: bool{true/false}
        SegmentSequenceSteps[]
          MatchType: string{"UNSPECIFIED_COHORT_TYPE","PRECEDES","IMMEDIATELY_PRECEDES"}
          OrFiltersForSegment[]:
            SegmentFilterClauses[]:
              DimensionFilter:
                CaseSensitive: bool{true/false}
                DimensionName: string{}
                Expressions[]: string{}
                MaxComparisonValue: string{}
                MinComparisonValue: string{}
                Operator: string{"OPERATOR_UNSPECIFIED","REGEXP","BEGINS_WITH","ENDS_WITH","PARTIAL","EXACT","IN_LIST","NUMERIC_LESS_THAN","NUMERIC_GREATER_THAN","NUMERIC_BETWEEN"}
                ForceSendFields[] string{}
              MetricFilter:
                ComparisonValue: string{}
                MaxComparisonValue: string{}
                MetricName: string{}
                Operator: string{"LESS_THAN","GREATER_THAN","EQUAL","BETWEEN"} // LT,GT?
                Scope: string{"UNSPECIFIED_SCOPE","PRODUCT","HIT","SESSION","USER"}
                ForceSendFields[] string{}
              Not: bool{true/false}
              ForceSendFields[] string{}
            ForceSendFields[] string{}
        ForceSendFields[] string{}
      SimpleSegment:
        OrFiltersForSegment[]
          SegmentFilterClauses[]:
            DimensionFilter:
              CaseSensitive: bool{true/false}
              DimensionName: string{}
              Expressions[]: string{}
              MaxComparisonValue: string{}
              MinComparisonValue: string{}
              Operator: string{"OPERATOR_UNSPECIFIED","REGEXP","BEGINS_WITH","ENDS_WITH","PARTIAL","EXACT","IN_LIST","NUMERIC_LESS_THAN","NUMERIC_GREATER_THAN","NUMERIC_BETWEEN"}
              ForceSendFields[] string{}
            MetricFilter:
              ComparisonValue: string{}
              MaxComparisonValue: string{}
              MetricName: string{}
              Operator: string{"LESS_THAN","GREATER_THAN","EQUAL","BETWEEN"} // LT,GT?
              Scope: string{"UNSPECIFIED_SCOPE","PRODUCT","HIT","SESSION","USER"}
              ForceSendFields[] string{}
            Not: bool{true/false}
            ForceSendFields[] string{}
          ForceSendFields[] string{}
        ForceSendFields[] string{}
      ForceSendFields[] string{}
  UserSegment:
    SegmentFilters[]:
      Not: bool{true/false}
      SequenceSegment:
        FirstStepShouldMatchFirstHit: bool{true/false}
        SegmentSequenceSteps[]
          MatchType: string{"UNSPECIFIED_COHORT_TYPE","PRECEDES","IMMEDIATELY_PRECEDES"}
          OrFiltersForSegment[]:
            SegmentFilterClauses[]:
              DimensionFilter:
                CaseSensitive: bool{true/false}
                DimensionName: string{}
                Expressions[]: string{}
                MaxComparisonValue: string{}
                MinComparisonValue: string{}
                Operator: string{"OPERATOR_UNSPECIFIED","REGEXP","BEGINS_WITH","ENDS_WITH","PARTIAL","EXACT","IN_LIST","NUMERIC_LESS_THAN","NUMERIC_GREATER_THAN","NUMERIC_BETWEEN"}
                ForceSendFields[] string{}
              MetricFilter:
                ComparisonValue: string{}
                MaxComparisonValue: string{}
                MetricName: string{}
                Operator: string{"LESS_THAN","GREATER_THAN","EQUAL","BETWEEN"} // LT,GT?
                Scope: string{"UNSPECIFIED_SCOPE","PRODUCT","HIT","SESSION","USER"}
                ForceSendFields[] string{}
              Not: bool{true/false}
              ForceSendFields[] string{}
            ForceSendFields[] string{}
        ForceSendFields[] string{}
      SimpleSegment:
        OrFiltersForSegment[]
          SegmentFilterClauses[]:
            DimensionFilter:
              CaseSensitive: bool{true/false}
              DimensionName: string{}
              Expressions[]: string{}
              MaxComparisonValue: string{}
              MinComparisonValue: string{}
              Operator: string{"OPERATOR_UNSPECIFIED","REGEXP","BEGINS_WITH","ENDS_WITH","PARTIAL","EXACT","IN_LIST","NUMERIC_LESS_THAN","NUMERIC_GREATER_THAN","NUMERIC_BETWEEN"}
              ForceSendFields[] string{}
            MetricFilter:
              ComparisonValue: string{}
              MaxComparisonValue: string{}
              MetricName: string{}
              Operator: string{"LESS_THAN","GREATER_THAN","EQUAL","BETWEEN"} // LT,GT?
              Scope: string{"UNSPECIFIED_SCOPE","PRODUCT","HIT","SESSION","USER"}
              ForceSendFields[] string{}
            Not: bool{true/false}
            ForceSendFields[] string{}
          ForceSendFields[] string{}
        ForceSendFields[] string{}
      ForceSendFields[] string{}
  ForceSendFields[] string{}
*/
