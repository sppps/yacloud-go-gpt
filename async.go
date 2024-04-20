package yacloud_gpt

import (
	"fmt"
	"log"
	"time"
)

type YandexGptRestAsync struct {
	FolderId               string
	ApiKey                 string
	IAMToken               string
	Logger                 *log.Logger
	BaseUrl                string
	OperationCheckInterval time.Duration
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
	if err != nil {
		return res, err
	}
	for !op.Done {
		time.Sleep(durationOrDefault(s.OperationCheckInterval, 1500*time.Millisecond))
		op, err = s.GetOperationResult(op.Id)
		if err != nil {
			return res, err
		}
	}
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
	res, err = callRestApi[asyncOperation](restApiCall{
		Endpoint: "completionAsync",
		ApiKey:   s.ApiKey,
		IAMToken: s.IAMToken,
		BaseUrl:  s.BaseUrl,
		Logger:   s.Logger,
		Params: completionReq{
			CompletionRequest: req,
			FolderId:          s.FolderId,
		},
	})
	return res, err
}

func (s YandexGptRestAsync) GetOperationResult(id string) (res asyncOperation, err error) {
	res, err = callRestApi[asyncOperation](restApiCall{
		BaseUrl:  "https://llm.api.cloud.yandex.net/operations",
		Endpoint: id,
		Method:   "GET",
		ApiKey:   s.ApiKey,
		IAMToken: s.IAMToken,
		Logger:   s.Logger,
	})
	return res, err
}
