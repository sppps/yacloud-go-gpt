package yacloud_gpt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type YandexGptRestAsync struct {
	FolderId               string
	ApiKey                 string
	IAMToken               string
	BaseUrl                string
	Logger                 *log.Logger
	OperationCheckInterval time.Duration
	restApiCaller
}

type asyncOperation struct {
	Id          string             `json:"id"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"createdAt"`
	CreatedBy   string             `json:"createdBy"`
	ModifiedAt  time.Time          `json:"modifiedAt"`
	Done        bool               `json:"done"`
	Metadata    map[string]any     `json:"metadata"`
	Error       apiError           `json:"error"`
	Response    CompletionResponse `json:"response"`
}

func (s YandexGptRestAsync) Completion(req CompletionRequest) (res CompletionResponse, err error) {
	op, err := s.CompletionAsync(req)
	for !op.Done {
		time.Sleep(durationOrDefault(s.OperationCheckInterval, 1500*time.Millisecond))
		op, err = s.GetOperationResult(op.Id)
	}
	s.Logger.Println(op.Response)
	return res, err
}

func (s YandexGptRestAsync) formatModelUri(uri ModelUri) ModelUri {
	if uri == YandexGptPro || uri == YandexGptLite || uri == YandexGptSummarization {
		return ModelUri(fmt.Sprintf("gpt://%s/%s/latest", s.FolderId, uri))
	}
	return uri
}

func (s YandexGptRestAsync) CompletionAsync(req CompletionRequest) (res asyncOperation, err error) {
	req.ModelUri = s.formatModelUri(req.ModelUri)
	d, err := s.callRestApi("completionAsync", completionReq{
		CompletionRequest: req,
		FolderId:          s.FolderId,
	})
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(d, &res)
	return res, err
}

func (s YandexGptRestAsync) GetOperationResult(id string) (res asyncOperation, err error) {
	url := fmt.Sprintf("https://llm.api.cloud.yandex.net/operations/%s", id)
	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return res, err
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
		return res, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(b, &res)
	return res, err
}
