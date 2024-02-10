package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type Item struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

var (
	discord_Token = ""
	StopBot       = make(chan bool)
	vcsession     *discordgo.VoiceConnection
	qiita_Url     = "https://qiita.com/api/v2/items?per_page=10&query="
)

func main() {
	query := "kubernetes"
	qiita_api_result := qiita_api_query(query)

	for _, item := range qiita_api_result {
		fmt.Printf("%s %s\n", item.Title, item.Url)
	}

}

// qiita_apiでタイトルとURLを表示する
func qiita_api_query(query string) []Item {
	resp, err := http.Get(qiita_Url + query)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data []Item

	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	return data
}
