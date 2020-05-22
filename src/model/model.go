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
	var tempShortUrl ShortURL

	err := database.DbMap.SelectOne(&tempShortUrl, "SELECT * FROM urls WHERE original=?", shortUrl.Original)
	if err == sql.ErrNoRows {
		var urlCode string
		urlCodeLength := config.Config.BaseURLLength

		for {
			for i := 0; i < 10; i++ {
				urlCode = util.RandomString(urlCodeLength)
				ret, err := database.DbMap.SelectInt("SELECT COUNT(*) FROM urls WHERE code=?", urlCode)
				if err != nil {
					return false, nil
				}

				if ret == 0 {
					goto success
				}
			}

			urlCodeLength++
		}

	success:
		shortUrl.Code = urlCode
		return true, nil
	} else if err != nil {
		return false, err
	}

	shortUrl.Code = tempShortUrl.Code
	return false, nil
}

func InitModels() error {
	database.DbMap.AddTableWithName(ShortURL{}, "urls")

	err := database.DbMap.CreateTablesIfNotExists()
	if err != nil {
		return err
	}

	return nil
}
