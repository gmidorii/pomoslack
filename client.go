package pomoslack

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	SQLiteFile string
}

type SlackTmp struct {
	Now        string
	StartTimes []string
}

const message = `{{ .Now }} のポモドーロ
{{ range $i, $v := .StartTimes }}* {{ $v }}
{{ end }}`

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

	log.Println(s)
	return nil
}

func exec(r result, now time.Time) (string, error) {
	tmp := SlackTmp{
		Now:        now.Format(layout),
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
