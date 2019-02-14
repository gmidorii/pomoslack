package pomoslack

import (
	"errors"
	"fmt"
	"strings"

	sl "github.com/nlopes/slack"
)

type Notifier interface {
	Send(title, dest, body string, ok bool) error
}

type slackNotify struct {
	token    string
	okColor  string
	errColor string
	mentions string
}

func NewSlack(token, okColor, errColor string, mentions string) Notifier {
	var mentionStr string
	for _, m := range strings.Split(mentions, ",") {
		mentionStr = fmt.Sprintf("%v<%v>,", mentionStr, m)
	}

	return slackNotify{
		token:    token,
		okColor:  okColor,
		errColor: errColor,
		mentions: mentionStr,
	}
}

func (s slackNotify) Send(title, dest, body string, ok bool) error {
	if s.token == "" {
		return errors.New("failed send message: token is empty")
	}
	client := sl.New(s.token)
	slackBody := fmt.Sprintf("%v\n%v", s.mentions, body)
	at := sl.Attachment{
		Color: s.okColor,
		Title: title,
		Text:  slackBody,
	}

	if !ok {
		at.Color = s.errColor
	}

	_, _, err := client.PostMessage(dest, sl.MsgOptionAttachments(at))
	return err
}
