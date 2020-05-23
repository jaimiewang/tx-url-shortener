package model

import (
	"errors"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/util"
)

type ShortURL struct {
	Id        int64  `db:"id, primarykey, autoincrement"`
	Time      int64  `db:"time"`
	IPAddress string `db:"ip_addr, size:15"`
	Original  string `db:"original, size:255"`
	Code      string `db:"code, size:11"`
	Used      int64  `db:"used"`
}

func (shortUrl *ShortURL) GenerateCode() error {
	if shortUrl.Id != 0 {
		return errors.New("already stored in database")
	}

	var urlCode string

	for i := config.Config.BaseCodeLength; i <= 11; i++ {
		for j := 0; j < 3; j++ {
			urlCode = util.RandomString(i, util.AsciiLetters)
			ret, err := database.DbMap.SelectInt("SELECT COUNT(*) FROM urls WHERE code=?", urlCode)
			if err != nil {
				return err
			}

			if ret == 0 {
				goto success
			}
		}
	}

	return errors.New("code cannot be generated")
success:
	shortUrl.Code = urlCode
	return nil
}

type APIKey struct {
	Id    int64  `db:"id, primarykey, autoincrement"`
	Time  int64  `db:"time"`
	Token string `db:"token, size:20"`
}

const APIKeySize = 20

func (apiKey *APIKey) GenerateToken() error {
	if apiKey.Id != 0 {
		return errors.New("already stored in database")
	}

	var token string

	for j := 0; j < 3; j++ {
		token = util.RandomString(APIKeySize, util.AsciiLetters)
		ret, err := database.DbMap.SelectInt("SELECT COUNT(*) FROM api_keys WHERE token=?", token)
		if err != nil {
			return err
		}

		if ret == 0 {
			goto success
		}
	}

	return errors.New("token cannot be generated")
success:
	apiKey.Token = token
	return nil
}

func InitModels() error {
	database.DbMap.AddTableWithName(APIKey{}, "api_keys")
	database.DbMap.AddTableWithName(ShortURL{}, "urls")

	err := database.DbMap.CreateTablesIfNotExists()
	if err != nil {
		return err
	}

	return nil
}
