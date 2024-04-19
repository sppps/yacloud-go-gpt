package yacloud_gpt

import (
	"fmt"
)

type YandexGptRest struct {
	restApiCaller
}

const defaultBaseUrl = "https://llm.api.cloud.yandex.net/foundationModels/v1"

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type completionReq struct {
	CompletionRequest
	FolderId string `json:"folder_id"`
}

func (s YandexGptRest) formatModelUri(uri ModelUri) ModelUri {
	if uri == YandexGptPro || uri == YandexGptLite || uri == YandexGptSummarization {
		return ModelUri(fmt.Sprintf("gpt://%s/%s/latest", s.FolderId, uri))
	}
	return uri
}

func (s YandexGptRest) Completion(req CompletionRequest) (res CompletionResponse, err error) {
	req.ModelUri = s.formatModelUri(req.ModelUri)
	_, err = s.callRestApi("completion", completionReq{
		CompletionRequest: req,
		FolderId:          s.FolderId,
	})
	return res, err
}
