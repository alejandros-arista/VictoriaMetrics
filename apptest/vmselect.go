package apptest

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
)

// Vmselect holds the state of a vmselect app and provides vmselect-specific
// functions.
type Vmselect struct {
	*app
	*ServesMetrics

	httpListenAddr          string
	clusternativeListenAddr string
	cli                     *Client
}

// StartVmselect starts an instance of vmselect with the given flags. It also
// sets the default flags and populates the app instance state with runtime
// values extracted from the application log (such as httpListenAddr)
func StartVmselect(instance string, flags []string, cli *Client) (*Vmselect, error) {
	app, stderrExtracts, err := startApp(instance, "../../bin/vmselect", flags, &appOptions{
		defaultFlags: map[string]string{
			"-httpListenAddr":          "127.0.0.1:0",
			"-clusternativeListenAddr": "127.0.0.1:0",
		},
		extractREs: []*regexp.Regexp{
			httpListenAddrRE,
			vmselectAddrRE,
		},
	})
	if err != nil {
		return nil, err
	}

	return &Vmselect{
		app: app,
		ServesMetrics: &ServesMetrics{
			metricsURL: fmt.Sprintf("http://%s/metrics", stderrExtracts[0]),
			cli:        cli,
		},
		httpListenAddr:          stderrExtracts[0],
		clusternativeListenAddr: stderrExtracts[1],
		cli:                     cli,
	}, nil
}

// ClusternativeListenAddr returns the address at which the vmselect process is
// listening for connections from other vmselect apps.
func (app *Vmselect) ClusternativeListenAddr() string {
	return app.clusternativeListenAddr
}

// PrometheusAPIV1Export is a test helper function that performs the export of
// raw samples in JSON line format by sending a HTTP POST request to
// /prometheus/api/v1/export vmselect endpoint.
//
// See https://docs.victoriametrics.com/url-examples/#apiv1export
func (app *Vmselect) PrometheusAPIV1Export(t *testing.T, query string, opts QueryOpts) *PrometheusAPIV1QueryResponse {
	t.Helper()

	exportURL := fmt.Sprintf("http://%s/select/%s/prometheus/api/v1/export", app.httpListenAddr, opts.getTenant())
	values := opts.asURLValues()
	values.Add("match[]", query)
	values.Add("format", "promapi")
	res := app.cli.PostForm(t, exportURL, values, getExpectedResponse(opts.ExpectedResponseCode, http.StatusOK))
	return NewPrometheusAPIV1QueryResponse(t, res)
}

// PrometheusAPIV1Query is a test helper function that performs PromQL/MetricsQL
// instant query by sending a HTTP POST request to /prometheus/api/v1/query
// vmselect endpoint.
//
// See https://docs.victoriametrics.com/url-examples/#apiv1query
func (app *Vmselect) PrometheusAPIV1Query(t *testing.T, query string, opts QueryOpts) *PrometheusAPIV1QueryResponse {
	t.Helper()

	queryURL := fmt.Sprintf("http://%s/select/%s/prometheus/api/v1/query", app.httpListenAddr, opts.getTenant())
	values := opts.asURLValues()
	values.Add("query", query)

	res := app.cli.PostForm(t, queryURL, values, getExpectedResponse(opts.ExpectedResponseCode, http.StatusOK))
	return NewPrometheusAPIV1QueryResponse(t, res)
}

// PrometheusAPIV1QueryRange is a test helper function that performs
// PromQL/MetricsQL range query by sending a HTTP POST request to
// /prometheus/api/v1/query_range vmselect endpoint.
//
// See https://docs.victoriametrics.com/url-examples/#apiv1query_range
func (app *Vmselect) PrometheusAPIV1QueryRange(t *testing.T, query string, opts QueryOpts) *PrometheusAPIV1QueryResponse {
	t.Helper()

	queryURL := fmt.Sprintf("http://%s/select/%s/prometheus/api/v1/query_range", app.httpListenAddr, opts.getTenant())
	values := opts.asURLValues()
	values.Add("query", query)

	res := app.cli.PostForm(t, queryURL, values, getExpectedResponse(opts.ExpectedResponseCode, http.StatusOK))
	return NewPrometheusAPIV1QueryResponse(t, res)
}

// PrometheusAPIV1Series sends a query to a /prometheus/api/v1/series endpoint
// and returns the list of time series that match the query.
//
// See https://docs.victoriametrics.com/url-examples/#apiv1series
func (app *Vmselect) PrometheusAPIV1Series(t *testing.T, matchQuery string, opts QueryOpts) *PrometheusAPIV1SeriesResponse {
	t.Helper()

	seriesURL := fmt.Sprintf("http://%s/select/%s/prometheus/api/v1/series", app.httpListenAddr, opts.getTenant())
	values := opts.asURLValues()
	values.Add("match[]", matchQuery)

	res := app.cli.PostForm(t, seriesURL, values, getExpectedResponse(opts.ExpectedResponseCode, http.StatusOK))
	return NewPrometheusAPIV1SeriesResponse(t, res)
}


// PrometheusAPIV1Labels sends a query to a /prometheus/api/v1/labels endpoint
// and returns the list of label names.
//
// See https://docs.victoriametrics.com/url-examples/#apiv1labels
func (app *Vmselect) PrometheusAPIV1Labels(t *testing.T, opts QueryOpts) *PrometheusAPIV1LabelsResponse {
	t.Helper()

	labelsURL := fmt.Sprintf("http://%s/select/%s/prometheus/api/v1/labels", app.httpListenAddr, opts.getTenant())
	values := opts.asURLValues()
	res := app.cli.PostForm(t, labelsURL, values, getExpectedResponse(opts.ExpectedResponseCode, http.StatusOK))
	return NewPrometheusAPIV1LabelsResponse(t, res)
}

// PrometheusAPIV1LabelValues sends a query to a /prometheus/api/v1/label/{label}/values endpoint
// and returns the list of label values for the specified label name.
//
// See https://docs.victoriametrics.com/url-examples/#apiv1labelvalues
func (app *Vmselect) PrometheusAPIV1LabelValues(t *testing.T, labelName string, opts QueryOpts) *PrometheusAPIV1LabelValuesResponse {
	t.Helper()

	valuesURL := fmt.Sprintf("http://%s/select/%s/prometheus/api/v1/label/%s/values", app.httpListenAddr, opts.getTenant(), labelName)
	values := opts.asURLValues()
	res := app.cli.PostForm(t, valuesURL, values, getExpectedResponse(opts.ExpectedResponseCode, http.StatusOK))
	return NewPrometheusAPIV1LabelValuesResponse(t, res)
}

// DeleteSeries sends a query to a /prometheus/api/v1/admin/tsdb/delete_series
//
// See https://docs.victoriametrics.com/url-examples/#apiv1admintsdbdelete_series
func (app *Vmselect) DeleteSeries(t *testing.T, matchQuery string, opts QueryOpts) {
	t.Helper()

	seriesURL := fmt.Sprintf("http://%s/delete/%s/prometheus/api/v1/admin/tsdb/delete_series", app.httpListenAddr, opts.getTenant())
	values := opts.asURLValues()
	values.Add("match[]", matchQuery)

	res := app.cli.PostForm(t, seriesURL, values, http.StatusNoContent)
	if res != "" {
		t.Fatalf("unexpected non-empty DeleteSeries response=%q", res)
	}
}

// String returns the string representation of the vmselect app state.
func (app *Vmselect) String() string {
	return fmt.Sprintf("{app: %s httpListenAddr: %q}", app.app, app.httpListenAddr)
}
