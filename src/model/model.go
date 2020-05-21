package model

import (
	"database/sql"
	"time"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/util"
)

type ShortURL struct {
	Id        int       `db:"id, primarykey, autoincrement"`
	Time      time.Time `db:"time"`
	IPAddress string    `db:"ip_addr"`
	Original  string    `db:"original"`
	Code      string    `db:"code"`
}

func (shortUrl *ShortURL) GenerateCode() (bool, error) {
	var urlCode string
	var tempShortUrl ShortURL

	err := database.DbMap.SelectOne(&tempShortUrl, "SELECT * FROM urls WHERE original=?", shortUrl.Original)
	if err == sql.ErrNoRows {
		urlCodeLength := config.Config.BaseURLLength
		urlsCount, err := database.DbMap.SelectInt("SELECT COUNT(*) FROM urls")
		if err != nil {
			return false, err
		}

		var counter int64

		for {
			counter = 0

			for {
				if urlsCount >= 4 && counter >= urlsCount/4 {
					break
				}

				urlCode = util.RandomString(urlCodeLength)
				ret, err := database.DbMap.SelectInt("SELECT COUNT(*) FROM urls WHERE code=?", urlCode)
				if err != nil {
					return false, nil
				}

				if ret == 0 {
					break
				}

				counter++
			}

			if urlsCount >= 1 && counter == urlsCount {
				urlCodeLength += 1
				continue
			} else {
				break
			}
		}

		shortUrl.Code = urlCode
		return true, nil
	} else if err != nil {
		return false, err
	}

	*shortUrl = tempShortUrl
	return false, nil
}
