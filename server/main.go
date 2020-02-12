package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yi-jiayu/nationstates"
)

type Notifier struct {
	PollInterval     time.Duration
	Client           *nationstates.Client
	Nation           string
	AdditionalShards []string
	Callback         func(notice nationstates.Notice, nation nationstates.Nation)

	ticker     *time.Ticker
	lastOffset int
}

type SendMessageRequest struct {
	ChatID      int    `json:"chat_id"`
	Text        string `json:"text"`
	ParseMode   string `json:"parse_mode"`
	ReplyMarkup string `json:"reply_markup"`
}

type InlineKeybardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url"`
	CallbackData string `json:"callback_data"`
}

type TelegramResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

func (n Notifier) poll() {
	log.Println("polling for notices")
	nation, err := n.Client.GetNation(n.Nation, append(n.AdditionalShards, "notices"), map[string]interface{}{"from": n.lastOffset})
	if err != nil {
		return
	}
	if notices := nation.Notices; len(notices) > 0 {
		log.Printf("got %d new notices\n", len(notices))
		n.lastOffset = notices[0].Timestamp
		for i := 0; i < len(notices); i++ {
			n.Callback(notices[len(notices)-i-1], nation)
		}
	}
}
func (n *Notifier) Start() {
	if n.ticker != nil {
		n.ticker.Stop()
	}
	n.poll()
	n.ticker = time.NewTicker(n.PollInterval)
	for range n.ticker.C {
		n.poll()
	}
}

func (r SendMessageRequest) Do(token string) error {
	u := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(r)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, u, &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var tgRes TelegramResponse
	err = json.NewDecoder(res.Body).Decode(&tgRes)
	if err != nil {
		return err
	}
	if !tgRes.Ok {
		return errors.New(tgRes.Description)
	}
	return nil
}

func sendMessage(token string, chatID int, text string) error {
	return SendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "Markdown",
	}.Do(token)
}

func sendMessageWithInlineKeyboard(token string, chatID int, text string, buttons [][]InlineKeyboardButton) error {
	replyMarkup, err := json.Marshal(InlineKeybardMarkup{InlineKeyboard: buttons})
	if err != nil {
		return err
	}
	return SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: string(replyMarkup),
	}.Do(token)
}

func getIssueID(notice nationstates.Notice) int {
	id, _ := strconv.Atoi(notice.URL[strings.LastIndex(notice.URL, "=")+1:])
	return id
}

func indexOfIssueWithID(issues []nationstates.Issue, id int) int {
	for i, issue := range issues {
		if issue.ID == id {
			return i
		}
	}
	return -1
}

func sendIssue(token string, chatID int, notice nationstates.Notice, issues []nationstates.Issue) error {
	id := getIssueID(notice)
	index := indexOfIssueWithID(issues, id)
	if index < 0 {
		return nil
	}
	issue := issues[index]
	err := sendMessage(token, chatID, fmt.Sprintf("*New Issue: %s*\n%s", issue.Title, issue.Text))
	if err != nil {
		return err
	}
	for _, option := range issue.Options {
		err := sendMessageWithInlineKeyboard(token, chatID, option.Text, [][]InlineKeyboardButton{
			{
				InlineKeyboardButton{
					Text:         "Accept",
					CallbackData: fmt.Sprintf("answerIssue,%d,%d", issue.ID, option.ID),
				},
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func newCallback(token string, chatID int) func(notice nationstates.Notice, nation nationstates.Nation) {
	return func(notice nationstates.Notice, nation nationstates.Nation) {
		switch notice.Type {
		case nationstates.NoticeIssue:
			err := sendIssue(token, chatID, notice, nation.Issues)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

type Config struct {
	Autologin string `json:"autologin"`
	Token     string `json:"token"`
	ChatID    int    `json:"chat_id"`
	Nation    string `json:"nation"`
}

func getConfig() (Config, error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func main() {
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	client := &nationstates.Client{
		Autologin: config.Autologin,
	}
	notifier := Notifier{
		PollInterval:     time.Hour,
		Client:           client,
		Nation:           config.Nation,
		AdditionalShards: []string{"issues"},
		Callback:         newCallback(config.Token, config.ChatID),
	}
	go notifier.Start()
}
