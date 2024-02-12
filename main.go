package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type Item struct {
	// Title string `json:"title"`
	Url string `json:"url"`
}

var (
	StopBot   = make(chan os.Signal, 1)
	qiita_Url = "https://qiita.com/api/v2/items?per_page=10&query="
)

func main() {
	//環境変数の読み取り
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print(err)
	}
	discord_Token := os.Getenv("DISCORD_TOKEN")

	//qiita_apiにリクエストを投げる時のタグ名
	query := "kubernetes"
	qiita_api_result := qiita_api_query(query)

	discord_bot(qiita_api_result[0], discord_Token)
}

// qiita_apiでタイトルとURLを表示する
func qiita_api_query(query string) []Item {
	var data []Item

	resp, err := http.Get(qiita_Url + query)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//json形式をGoのデータ構造に変換
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	return data
}

// discordに送信する
func discord_bot(i Item, token string) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Print(err)
	}

	discord.AddHandler(onMessageCreate)
	// Channelid := os.Getenv("CHANNEL_ID")
	// sendMessage(Channelid, "a")

	// sendMessage(m.ChannelD)
	err = discord.Open()
	if err != nil {
		fmt.Print(err)
	}
	// signal.Notify(StopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	defer discord.Close()
	<-StopBot
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// fmt.Printf("%20s %20s(%20s) > %s\n", m.ChannelID, m.Content)
	clientid := os.Getenv("CLIENT_ID")
	u := m.Author
	if u.ID != clientid {
		sendMessage(s, m.ChannelID, "a")
	}
}

func sendMessage(s *discordgo.Session, ChannelID string, msg string) {
	_, err := s.ChannelMessageSend(ChannelID, msg)
	fmt.Print(ChannelID)
	if err != nil {
		fmt.Print(err)
	}
}
