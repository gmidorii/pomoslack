package pomoslack

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type result struct {
	count      int
	startTimes []time.Time
}

type Entity struct {
	PK        string
	StartTime string
	EndTime   string
}

type client struct {
	db *sql.DB
}

func NewClient(db *sql.DB) client {
	return client{
		db: db,
	}
}

const query = `select Z_PK, strftime('%H:%M', ZENDEDTIME, 'unixepoch'), strftime('%H:%M', ZSTARTEDTIME, 'unixepoch') from ZTASK where ZENDEDTIME IS NOT NULL and strftime('%m-%d', ZSTARTEDTIME, 'unixepoch') = ?`

func (c *client) fetch(ctx context.Context, now time.Time) (result, error) {
	rows, err := c.db.QueryContext(ctx, query, now.Format(layoutDate))
	if err != nil {
		return result{}, err
	}
	defer rows.Close()

	var entities []Entity
	for rows.Next() {
		var e Entity
		if err := rows.Scan(&e.PK, &e.EndTime, &e.StartTime); err != nil {
			return result{}, err
		}
		entities = append(entities, e)
	}

	startTimes := make([]time.Time, len(entities))
	for i, e := range entities {
		st, err := convTime(layoutTime, e.StartTime, now)
		if err != nil {
			return result{}, err
		}
		startTimes[i] = st
	}

	return result{count: len(entities), startTimes: startTimes}, nil
}

func convTime(layout, tstr string, now time.Time) (time.Time, error) {
	t, err := time.Parse(layout, tstr)
	if err != nil {
		return time.Time{}, err
	}
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return time.Time{}, err
	}

	t = t.In(time.UTC).In(jst)

	return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, jst), nil
}
