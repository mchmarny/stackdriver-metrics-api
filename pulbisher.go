package main

import (
	"context"
	"time"

	"github.com/pkg/errors"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

const (
	customerMetricNamespace = "custom.googleapis.com/workflow/metric"
)

func post(ctx context.Context, metric string, when time.Time, count int64) error {

	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to create monitor client")
	}

	dataPoint := &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			StartTime: &googlepb.Timestamp{Seconds: when.Unix()},
			EndTime:   &googlepb.Timestamp{Seconds: when.Unix()},
		},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_Int64Value{Int64Value: m},
		},
	}

	tsRequest := &monitoringpb.CreateTimeSeriesRequest{
		Name: monitoring.MetricProjectPath(projectID),
		TimeSeries: []*monitoringpb.TimeSeries{
			{
				Metric: &metricpb.Metric{
					Type:   customerMetricNamespace,
					Labels: map[string]string{"metric_name": metric},
				},
				Resource: &monitoredrespb.MonitoredResource{
					Type:   "global",
					Labels: map[string]string{"project_id": projectID},
				},
				Points: []*monitoringpb.Point{dataPoint},
			},
		},
	}

	if err := client.CreateTimeSeries(ctx, tsRequest); err != nil {
		return errors.Wrap(err, "Error writting time series data")
	}

	return nil

}
