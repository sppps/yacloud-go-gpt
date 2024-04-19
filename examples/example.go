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
				Text: "Ты — переводчик и маркетолог. Напиши описание товара для маркетплейса на русском языке. Используй заданное оригинальное описание товара.",
			},
			{
				Role: yacloud_gpt.ModeUser,
				Text: "<div id=\"offer-template-0\"></div><p><span style=\"font-size: 24.0pt;font-family: simhei;\"><img alt=\"i2mSj/ukz9QsY8Ty6eyZZboAdJICZLpP7Cge\" alt=\"Свеча зажигания\" src=\"https://cbu01.alicdn.com/img/ibank/2015/353/562/2124265353_912732929.jpg\"><br><br><img alt=\"TIM Picture 20180814150209\" src=\"https://cbu01.alicdn.com/img/ibank/2018/377/248/9248842773_912732929.jpg\"><br><br><img alt=\"TIM Picture 20180814150119\" src=\"https://cbu01.alicdn.com/img/ibank/2018/536/158/9248851635_912732929.jpg\"><br><br>Угол Комнаты Компании</span><br><img alt=\"QQ Picture 20170708152625\" height=\"592.5\" src=\"https://cbu01.alicdn.com/img/ibank/2017/690/378/4413873096_912732929.jpg\" width=\"790\"><br><br><img alt=\"QQ Picture 20170708152544\" height=\"592.5\" src=\"https://cbu01.alicdn.com/img/ibank/2017/858/922/4408229858_912732929.jpg\" width=\"790\"><br><br><img alt=\"QQ Picture 20170708152511\" height=\"1053.3333333333335\" src=\"https://cbu01.alicdn.com/img/ibank/2017/365/532/4408235563_912732929.jpg\" width=\"790\"><br><br></p>",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
