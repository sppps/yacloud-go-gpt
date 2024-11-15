package main

import (
	"fmt"
	"log"
	"os"

	yacloud_gpt "github.com/sppps/yacloud-go-gpt"
)

func main() {
	gpt := yacloud_gpt.YandexGptRestAsync{
		FolderId: os.Getenv("YACLOUD_GPT_FOLDER_ID"),
		ApiKey:   os.Getenv("YACLOUD_GPT_API_KEY"),
		Logger:   log.Default(),
	}
	resp, err := gpt.Completion(yacloud_gpt.CompletionRequest{
		ModelUri: yacloud_gpt.YandexGptPro,
		Messages: []yacloud_gpt.CompletionMessage{
			{
				Role: yacloud_gpt.ModeSystem,
				Text: "You are an AI programming assistant. Follow the user's requirements carefully and to the letter." +
					"First, think step-by-step and describe your plan for what to build in pseudocode, written out in great detail. " +
					"Then, output the code in a single code block. Minimize any other prose.",
			},
			{
				Role: yacloud_gpt.ModeUser,
				Text: "hello_world.go",
			},
		},
		CompletionOptions: &yacloud_gpt.CompletionOptions{
			Temperature: 0.75,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
