package model

import (
	"database/sql"
	"errors"
	"time"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/util"
)

type ShortURL struct {
	Id        int64     `db:"id, primarykey, autoincrement"`
	Time      time.Time `db:"time"`
	IPAddress string    `db:"ip_addr,size:15"`
	Original  string    `db:"original,size:255"`
	Code      string    `db:"code,size:255"`
}

func (shortUrl *ShortURL) GenerateCode() (bool, error) {
	var tempShortUrl ShortURL

	err := database.DbMap.SelectOne(&tempShortUrl, "SELECT * FROM urls WHERE original=?", shortUrl.Original)
	if err == sql.ErrNoRows {
		var urlCode string

		for i := config.Config.BaseURLLength; i <= 255; i++ {
			for j := 0; j < 3; j++ {
				urlCode = util.RandomString(i)
				ret, err := database.DbMap.SelectInt("SELECT COUNT(*) FROM urls WHERE code=?", urlCode)
				if err != nil {
					return false, err
				}

				if ret == 0 {
					goto success
				}
			}
		}

		return false, errors.New("code cannot be generated")
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
