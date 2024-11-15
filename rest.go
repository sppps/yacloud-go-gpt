package yacloud_gpt

import (
	"context"
	"fmt"
	"log"
)

type YandexGptRest struct {
	FolderId string
	ApiKey   string
	IAMToken string
	Logger   *log.Logger
	BaseUrl  string
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type completionReq struct {
	CompletionRequest
	FolderId string `json:"folder_id"`
}

func (s YandexGptRest) Completion(req CompletionRequest) (res CompletionResponse, err error) {
	return s.CompletionWithContext(context.Background(), req)
}

func (s YandexGptRest) CompletionWithContext(ctx context.Context, req CompletionRequest) (res CompletionResponse, err error) {
	req.ModelUri = fmt.Sprintf("gpt://%s/%s", s.FolderId, req.ModelUri)
	res, err = callRestApi[CompletionResponse](ctx, restApiCall{
		Endpoint: "completion",
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
