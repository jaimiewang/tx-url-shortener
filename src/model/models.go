package model

import "time"

type ShortURL struct {
	Id        int       `db:"id, primarykey, autoincrement"`
	Time      time.Time `db:"time"`
	IPAddress string    `db:"ip_addr"`
	Original  string    `db:"original"`
	Code      string    `db:"code"`
}
