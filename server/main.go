package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

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
		data, err := json.Marshal(CallbackData{
			Action:   "answerIssue",
			IssueID:  issue.ID,
			OptionID: option.ID,
		})
		if err != nil {
			return err
		}
		err = sendMessageWithInlineKeyboard(token, chatID, option.Text, [][]InlineKeyboardButton{
			{
				InlineKeyboardButton{
					Text:         "Accept",
					CallbackData: string(data),
				},
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func answerCallbackQuery(token, id string) error {
	u := fmt.Sprintf("https://api.telegram.org/bot%s/answerCallbackQuery", token)
	params := url.Values{}
	params.Add("callback_query_id", id)
	res, err := http.PostForm(u, params)
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

func newCallback(token string, chatID int) func(notice nationstates.Notice, nation nationstates.Nation) {
	return func(notice nationstates.Notice, nation nationstates.Nation) {
		switch notice.Type {
		case nationstates.NoticeIssue:
			err := sendIssue(token, chatID, notice, nation.Issues)
			if err != nil {
				log.Println(err)
			}
		case nationstates.NoticeRank:
			text := fmt.Sprintf("%s %s", notice.Who, notice.Text)
			err := sendMessage(token, chatID, text)
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

type Update struct {
	CallbackQuery *CallbackQuery `json:"callback_query"`
}

type CallbackQuery struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type CallbackData struct {
	Action   string `json:"a"`
	IssueID  int    `json:"iid"`
	OptionID int    `json:"oid"`
}

func newUpdateHandler(client *nationstates.Client, nation, token string, chatID int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var u Update
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			return
		}
		callbackQuery := u.CallbackQuery
		if callbackQuery == nil {
			return
		}
		var d CallbackData
		err = json.Unmarshal([]byte(callbackQuery.Data), &d)
		if err != nil {
			return
		}
		switch d.Action {
		case "answerIssue":
			conseq, err := client.AnswerIssue(nation, d.IssueID, d.OptionID)
			if err != nil {
				log.Println(err)
				return
			}
			var text string
			if conseq.Error != "" {
				text = conseq.Error
			} else {
				talkingPoint := []rune(conseq.Desc)
				talkingPoint[0] = unicode.ToUpper(talkingPoint[0])
				headlines := strings.Join(conseq.Headlines, "\n")
				text = fmt.Sprintf("*The Talking Point*\n%s.\n\n*Recent Headlines*\n%s", string(talkingPoint), headlines)
			}
			err = sendMessage(token, chatID, text)
			if err != nil {
				log.Println(err)
			}
			err = answerCallbackQuery(token, callbackQuery.ID)
			if err != nil {
				log.Println(err)
			}
		}
	}
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
	http.ListenAndServe(":8080", http.HandlerFunc(newUpdateHandler(client, config.Nation, config.Token, config.ChatID)))
}
