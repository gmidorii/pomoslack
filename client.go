package pomoslack

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	SQLiteFile string
	Slack      Slack
}

type Slack struct {
	Token   string
	Channel string
}

type SlackTmp struct {
	Count      int
	StartTimes []string
}

const (
	title   = "%v のポモドーロ"
	message = `count: {{ .Count }}
{{ range $i, $v := .StartTimes }}* {{ $v }}
{{ end }}`
)

const (
	layout     = "2006-01-02"
	layoutTime = "15:04"
	layoutDate = "01-02"
)

func Run(c Config) error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%v", c.SQLiteFile))
	if err != nil {
		return errors.Wrap(err, "failed sqlite open")
	}
	defer db.Close()

	cli := NewClient(db)
	now := time.Now()
	r, err := cli.fetch(context.Background(), now)
	if err != nil {
		return errors.Wrap(err, "failed fetch from sqlite")
	}

	s, err := exec(r, now)
	if err != nil {
		return errors.Wrap(err, "failed create tamplate")
	}

	slack := NewSlack(c.Slack.Token, "#006400", "#dc143c", "")
	if err := slack.Send(fmt.Sprintf(title, now.Format(layoutDate)), c.Slack.Channel, s, true); err != nil {
		return errors.Wrap(err, "failed send slack")
	}
	return nil
}

func exec(r result, now time.Time) (string, error) {
	tmp := SlackTmp{
		Count:      r.count,
		StartTimes: convStrs(r.startTimes),
	}

	t, err := template.New("slack").Parse(message)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := t.Execute(&b, tmp); err != nil {
		return "", err
	}

	return b.String(), nil
}

func convStrs(ds []time.Time) []string {
	converted := make([]string, len(ds))
	for i, v := range ds {
		converted[i] = v.Format(layoutTime)
	}
	return converted
}
