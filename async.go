package yacloud_gpt

import (
	"context"
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
	return s.CompletionWithContext(context.Background(), req)
}

func (s YandexGptRestAsync) CompletionWithContext(ctx context.Context, req CompletionRequest) (res CompletionResponse, err error) {
	op, err := s.completionAsync(ctx, req)
	if err != nil {
		return res, err
	}
	for !op.Done && ctx.Err() == nil {
		time.Sleep(durationOrDefault(s.OperationCheckInterval, 1500*time.Millisecond))
		op, err = s.getOperationResult(ctx, op.Id)
		if err != nil {
			return res, err
		}
	}
	return op.Response, err
}

func (s YandexGptRestAsync) completionAsync(ctx context.Context, req CompletionRequest) (res asyncOperation, err error) {
	req.ModelUri = fmt.Sprintf("gpt://%s/%s", s.FolderId, req.ModelUri)
	res, err = callRestApi[asyncOperation](ctx, restApiCall{
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

func (s YandexGptRestAsync) getOperationResult(ctx context.Context, id string) (res asyncOperation, err error) {
	res, err = callRestApi[asyncOperation](ctx, restApiCall{
		BaseUrl:  "https://llm.api.cloud.yandex.net/operations",
		Endpoint: id,
		Method:   "GET",
		ApiKey:   s.ApiKey,
		IAMToken: s.IAMToken,
		Logger:   s.Logger,
	})
	return res, err
}
