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

type restApiCaller struct {
	BaseUrl  string
	Logger   *log.Logger
	FolderId string
	ApiKey   string
	IAMToken string
}

func (s restApiCaller) callRestApi(method string, params any) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", stringOrDefault(s.BaseUrl, defaultBaseUrl), method)

	if s.Logger != nil {
		s.Logger.Printf("yacloud translate: %s", url)
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if s.Logger != nil {
		s.Logger.Printf("yacloud translate: %s", string(body))
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")
	if len(s.ApiKey) > 0 {
		req.Header.Set("authorization", fmt.Sprintf("Api-Key %s", s.ApiKey))
	} else if len(s.IAMToken) > 0 {
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", s.IAMToken))
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if s.Logger != nil {
		s.Logger.Printf("yacloud translate: %s", string(d))
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
