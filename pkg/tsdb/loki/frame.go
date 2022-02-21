package loki

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// the generic prom&loki framing code creates dataframes for us,
// we need to adjust them a little to be like they are needed
// by the loki datasource
func adjustFrame(frame *data.Frame, query *lokiQuery) *data.Frame {
	labels := getFrameLabels(frame)

	// FIXME: is this unique to numeric stuff? what will the loki-streams get?
	isMetricFrame := frame.Meta.Type == data.FrameTypeTimeSeriesMany

	// FIXME: when instant-queries arrive, switch here based on QueryType
	isRangeQuery := true

	isMetricRanged := isMetricFrame && isRangeQuery

	name := formatName(labels, query)
	frame.Name = name

	if frame.Meta == nil {
		frame.Meta = &data.FrameMeta{}
	}

	if isMetricRanged {
		frame.Meta.ExecutedQueryString = "Expr: " + query.Expr + "\n" + "Step: " + query.Step.String()
	} else {
		frame.Meta.ExecutedQueryString = "Expr: " + query.Expr
	}

	timeFields, nonTimeFields := partitionFields(frame)

	for _, field := range timeFields {
		field.Name = "time"

		if isMetricRanged {
			if field.Config == nil {
				field.Config = &data.FieldConfig{}
			}
			field.Config.Interval = float64(query.Step.Milliseconds())
		}
	}

	for _, field := range nonTimeFields {
		field.Name = "value"
		if field.Config == nil {
			field.Config = &data.FieldConfig{}
		}
		field.Config.DisplayNameFromDS = name
	}

	// FIXME: i clear the frame-type here, so that the unit-test-snapshots do not need to be changed yet.
	// when this is merged, i will remove this and update the snaphsots in a separate PR.
	frame.Meta.Type = data.FrameTypeUnknown

	return frame
}

func formatNamePrometheusStyle(labels map[string]string) string {
	metricName := labels["__name__"]

	var parts []string

	for k, v := range labels {
		if k != "__name__" {
			parts = append(parts, fmt.Sprintf("%s=%q", k, v))
		}
	}

	sort.Strings(parts)

	return fmt.Sprintf("%s{%s}", metricName, strings.Join(parts, ", "))
}

//If legend (using of name or pattern instead of time series name) is used, use that name/pattern for formatting
func formatName(labels map[string]string, query *lokiQuery) string {
	if query.LegendFormat == "" {
		return formatNamePrometheusStyle(labels)
	}

	result := legendFormat.ReplaceAllFunc([]byte(query.LegendFormat), func(in []byte) []byte {
		labelName := strings.Replace(string(in), "{{", "", 1)
		labelName = strings.Replace(labelName, "}}", "", 1)
		labelName = strings.TrimSpace(labelName)
		if val, exists := labels[labelName]; exists {
			return []byte(val)
		}
		return []byte{}
	})

	return string(result)
}

func getFrameLabels(frame *data.Frame) map[string]string {
	labels := make(map[string]string)

	for _, field := range frame.Fields {
		for k, v := range field.Labels {
			labels[k] = v
		}
	}

	return labels
}

func partitionFields(frame *data.Frame) ([]*data.Field, []*data.Field) {
	var timeFields []*data.Field
	var nonTimeFields []*data.Field

	for _, field := range frame.Fields {
		if field.Type() == data.FieldTypeTime {
			timeFields = append(timeFields, field)
		} else {
			nonTimeFields = append(nonTimeFields, field)
		}
	}

	return timeFields, nonTimeFields
}
