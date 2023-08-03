package app

import (
	"fmt"
	"strconv"

	"github.com/mattermost/mattermost-perf-stats-cli/prometheus"
)

type DBEntry struct {
	Method    string
	TotalTime float64
	Count     float64
	Average   float64
}

type APIEntry struct {
	Handler   string
	TotalTime float64
	Count     float64
	Average   float64
}

type App struct {
	client *prometheus.Client
}

func New(endpoint string) *App {
	client := prometheus.New(endpoint)
	return &App{
		client: client,
	}
}

func (a *App) GetDBMetrics(timeRange string) (map[string]*DBEntry, error) {
	data := map[string]*DBEntry{}
	totalTimeMetrics, err := a.client.Query(fmt.Sprintf("sum(increase(mattermost_db_store_time_sum[%s]) and increase(mattermost_db_store_time_count[%s]) > 0) by (method)", timeRange, timeRange))
	if err != nil {
		return nil, err
	}
	for _, r := range totalTimeMetrics.Data.Result {
		calls, err := strconv.ParseFloat(r.Value[1].(string), 64)
		if err != nil {
			return nil, err
		}
		data[r.Metric["method"]] = &DBEntry{TotalTime: calls, Method: r.Metric["method"]}
	}

	callsMetrics, err := a.client.Query(fmt.Sprintf("sum(increase(mattermost_db_store_time_count[%s]) > 0) by (method)", timeRange))
	if err != nil {
		return nil, err
	}
	for _, r := range callsMetrics.Data.Result {
		entry := data[r.Metric["method"]]
		count, err := strconv.ParseFloat(r.Value[1].(string), 64)
		if err != nil {
			return nil, err
		}
		entry.Count = count
	}
	averageMetrics, err := a.client.Query(fmt.Sprintf("(sum(increase(mattermost_db_store_time_sum[%s])) by (method) / sum(increase(mattermost_db_store_time_count[%s]) > 0) by (method))", timeRange, timeRange))
	if err != nil {
		return nil, err
	}
	for _, r := range averageMetrics.Data.Result {
		entry := data[r.Metric["method"]]
		average, err := strconv.ParseFloat(r.Value[1].(string), 64)
		if err != nil {
			return nil, err
		}
		entry.Average = average
	}
	return data, nil
}

func (a *App) GetAPIMetrics(timeRange string) (map[string]*APIEntry, error) {
	data := map[string]*APIEntry{}
	totalTimeMetrics, err := a.client.Query(fmt.Sprintf("sum(increase(mattermost_api_time_sum[%s]) and increase(mattermost_api_time_count[%s]) > 0) by (handler)", timeRange, timeRange))
	if err != nil {
		return nil, err
	}
	for _, r := range totalTimeMetrics.Data.Result {
		calls, err := strconv.ParseFloat(r.Value[1].(string), 64)
		if err != nil {
			return nil, err
		}
		data[r.Metric["handler"]] = &APIEntry{TotalTime: calls, Handler: r.Metric["handler"]}
	}

	callsMetrics, err := a.client.Query(fmt.Sprintf("sum(increase(mattermost_api_time_count[%s]) > 0) by (handler)", timeRange))
	if err != nil {
		return nil, err
	}
	for _, r := range callsMetrics.Data.Result {
		entry := data[r.Metric["handler"]]
		count, err := strconv.ParseFloat(r.Value[1].(string), 64)
		if err != nil {
			return nil, err
		}
		entry.Count = count
	}
	averageMetrics, err := a.client.Query(fmt.Sprintf("(sum(increase(mattermost_api_time_sum[%s])) by (handler) / sum(increase(mattermost_api_time_count[%s]) > 0) by (handler))", timeRange, timeRange))
	if err != nil {
		return nil, err
	}
	for _, r := range averageMetrics.Data.Result {
		entry := data[r.Metric["handler"]]
		average, err := strconv.ParseFloat(r.Value[1].(string), 64)
		if err != nil {
			return nil, err
		}
		entry.Average = average
	}
	return data, nil
}
