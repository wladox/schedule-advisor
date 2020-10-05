package forecast

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wladox/schedule-advisor/pkg/network"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ForecastsResponse struct {
	IngestionRate []float64 `json:"ingestion_rate"`
	OperatorState []float64 `json:"operator_state"`
	Bandwidth     []float64 `json:"bandwidth"`
}

// Client to make the request to forecast service
type ForecastClient struct {
	doer network.HTTPDoer
	host string
}

// NewClient create new client
func NewClient(doer network.HTTPDoer, host string) *ForecastClient {
	return &ForecastClient{
		doer: doer,
		host: host,
	}
}

func (f *ForecastClient) GetForecasts(ctx context.Context, time int) (ForecastsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/", f.host), nil)
	if err != nil {
		return ForecastsResponse{}, err
	}

	q := req.URL.Query()
	q.Add("time", strconv.Itoa(time))
	req.URL.RawQuery = q.Encode()

	req = req.WithContext(ctx)
	res, err := f.doer.Do(req)
	if err != nil {
		return ForecastsResponse{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ForecastsResponse{}, err
	}

	var forecasts ForecastsResponse
	err = json.Unmarshal(body, &forecasts)
	if err != nil {
		return ForecastsResponse{}, err
	}

	return forecasts, nil
}
