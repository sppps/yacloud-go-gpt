package yacloud_gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func stringOrDefault(s1, s2 string) string {
	if len(s1) > 0 {
		return s1
	}
	return s2
}

func durationOrDefault(d1, d2 time.Duration) time.Duration {
	if d1 > 0 {
		return d1
	}
	return d2
}

type restApiCall struct {
	Method   string
	Params   any
	BaseUrl  string
	Logger   *log.Logger
	ApiKey   string
	IAMToken string
}

func callRestApi(req restApiCall) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", stringOrDefault(req.BaseUrl, defaultBaseUrl), req.Method)

	if req.Logger != nil {
		req.Logger.Printf("yacloud translate: %s", url)
	}

	body, err := json.Marshal(req.Params)
	if err != nil {
		return nil, err
	}

	if req.Logger != nil {
		req.Logger.Printf("yacloud translate: %s", string(body))
	}

	rreq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	rreq.Header.Set("content-type", "application/json")
	if len(req.ApiKey) > 0 {
		rreq.Header.Set("authorization", fmt.Sprintf("Api-Key %s", req.ApiKey))
	} else if len(req.IAMToken) > 0 {
		rreq.Header.Set("authorization", fmt.Sprintf("Bearer %s", req.IAMToken))
	}

	client := http.Client{}
	resp, err := client.Do(rreq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if req.Logger != nil {
		req.Logger.Printf("yacloud translate: %s", string(d))
	}

	if resp.StatusCode != http.StatusOK {
		var apiError apiError
		err = json.Unmarshal(d, &apiError)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("api error %d: %s", apiError.Code, apiError.Message)
	}

	return d, err
}
