package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/kiali/kiali/prometheus"
)

func extractServiceMetricsQuery(r *http.Request, namespace, service string) (*prometheus.ServiceMetricsQuery, error) {
	q := prometheus.ServiceMetricsQuery{
		Namespace: namespace,
		Service:   service}
	err := extractMetricsQueryParams(r, &q.MetricsQuery)
	return &q, err
}

func extractNamespaceMetricsQuery(r *http.Request, namespace, servicePattern string) (*prometheus.NamespaceMetricsQuery, error) {
	q := prometheus.NamespaceMetricsQuery{
		Namespace:      namespace,
		ServicePattern: servicePattern}
	err := extractMetricsQueryParams(r, &q.MetricsQuery)
	return &q, err
}

func extractMetricsQueryParams(r *http.Request, q *prometheus.MetricsQuery) error {
	q.FillDefaults()
	queryParams := r.URL.Query()
	if rateIntervals, ok := queryParams["rateInterval"]; ok && len(rateIntervals) > 0 {
		// Only first is taken into consideration
		q.RateInterval = rateIntervals[0]
	}
	if rateFuncs, ok := queryParams["rateFunc"]; ok && len(rateFuncs) > 0 {
		// Only first is taken into consideration
		if rateFuncs[0] != "rate" && rateFuncs[0] != "irate" {
			// Bad request
			return errors.New("Bad request, query parameter 'rateFunc' must be either 'rate' or 'irate'")
		}
		q.RateFunc = rateFuncs[0]
	}
	if durations, ok := queryParams["duration"]; ok && len(durations) > 0 {
		if num, err := strconv.Atoi(durations[0]); err == nil {
			q.Duration = time.Duration(num) * time.Second
		} else {
			// Bad request
			return errors.New("Bad request, cannot parse query parameter 'duration'")
		}
	}
	if steps, ok := queryParams["step"]; ok && len(steps) > 0 {
		if num, err := strconv.Atoi(steps[0]); err == nil {
			q.Step = time.Duration(num) * time.Second
		} else {
			// Bad request
			return errors.New("Bad request, cannot parse query parameter 'step'")
		}
	}
	if versions, ok := queryParams["version"]; ok && len(versions) > 0 {
		q.Version = versions[0]
	}
	if filters, ok := queryParams["filters[]"]; ok && len(filters) > 0 {
		q.Filters = filters
	}
	if lblsin, ok := queryParams["byLabelsIn[]"]; ok && len(lblsin) > 0 {
		q.ByLabelsIn = lblsin
	}
	if lblsout, ok := queryParams["byLabelsOut[]"]; ok && len(lblsout) > 0 {
		q.ByLabelsOut = lblsout
	}
	return nil
}
