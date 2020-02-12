package nationstates

import (
	"encoding/xml"
)

const (
	NoticeTelegram          = "TG"
	NoticeIssue             = "I"
	NoticeEndorsementGained = "END"
	NoticeEndorsementLost   = "UNEND"
	NoticeBanner            = "U"
	NoticeRank              = "T"
	NoticePolicy            = "P"
	NoticeTradingCards      = "C"
	NoticeRMBMention        = "RMB"
	NoticeRMBQuote          = "RMBQ"
	NoticeRMBLike           = "RMBL"
	NoticeDispatchMention   = "D"
	NoticeDispatchPin       = "DP"
	NoticeDispatchQuote     = "DQ"
	NoticeEmbassy           = "EMB"
	NoticeLoomingApocalypse = "X"
)

type Nation struct {
	XMLName      xml.Name     `xml:"NATION"`
	ID           string       `xml:"id,attr"`
	Consequences Consequences `xml:"ISSUE"`
	Issues       []Issue      `xml:"ISSUES>ISSUE"`
	Notices      []Notice     `xml:"NOTICES>NOTICE"`
}

type Issue struct {
	ID      int      `xml:"id,attr"`
	Title   string   `xml:"TITLE"`
	Text    string   `xml:"TEXT"`
	Options []Option `xml:"OPTION"`
}

type Consequences struct {
	Desc      string   `xml:"DESC"`
	Rankings  []Rank   `xml:"RANKINGS>RANK"`
	Headlines []string `xml:"HEADLINES>HEADLINE"`

	Error string `xml:"ERROR"`
}

type Rank struct {
	Score   float32 `xml:"SCORE"`
	Change  float32 `xml:"CHANGE"`
	PChange float32 `xml:"PCHANGE"`
}

type Option struct {
	ID   int    `xml:"id,attr"`
	Text string `xml:",chardata"`
}

type Notice struct {
	Text      string `xml:"TEXT"`
	Timestamp int    `xml:"TIMESTAMP"`
	Title     string `xml:"TITLE"`
	Who       string `xml:"WHO"`
	URL       string `xml:"URL"`
	Type      string `xml:"TYPE"`
}
