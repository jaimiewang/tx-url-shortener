package model

import (
	"errors"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/util"
)

var ErrAlreadyInDatabase = errors.New("already stored in database")
var ErrCannotBeGenerated = errors.New("cannot be generated")

type Model struct {
	ID int64 `db:"id, primarykey, autoincrement"`
}

type ShortURL struct {
	Model
	CreatedAt int64  `db:"created_at"`
	IPAddress string `db:"ip_addr, size:15"`
	Original  string `db:"original, size:255"`
	Code      string `db:"code, size:11"`
	Counter   int64  `db:"counter"`
}

func (shortURL *ShortURL) GenerateCode() error {
	if shortURL.ID != 0 {
		return ErrAlreadyInDatabase
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

	return ErrCannotBeGenerated
success:
	shortURL.Code = urlCode
	return nil
}

const APIKeySize = 20

type APIKey struct {
	Model
	CreatedAt int64  `db:"created_at"`
	Token     string `db:"token"`
}

func (apiKey *APIKey) GenerateToken() error {
	if apiKey.ID != 0 {
		return ErrAlreadyInDatabase
	}

	var token string
	var err error

	for j := 0; j < 3; j++ {
		token, err = util.RandomToken(APIKeySize)
		if err != nil {
			return err
		}

		ret, err := database.DbMap.SelectInt("SELECT COUNT(*) FROM api_keys WHERE token=?", token)
		if err != nil {
			return err
		}

		if ret == 0 {
			goto success
		}
	}

	return ErrCannotBeGenerated
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
