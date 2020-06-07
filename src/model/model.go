package model

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"gopkg.in/gorp.v2"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/util"
)

var ErrCannotBeGenerated = errors.New("model: cannot be generated")

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

func (shortURL *ShortURL) GenerateCode(trans *gorp.Transaction) error {
	var urlCode string

	for i := config.Config.BaseCodeLength; i <= 11; i++ {
		for j := 0; j < 3; j++ {
			urlCode = util.RandomString(i, util.AsciiLetters)
			ret, err := trans.SelectInt("SELECT COUNT(*) FROM urls WHERE code=?", urlCode)
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

func (shortURL *ShortURL) IsDoubled(trans *gorp.Transaction) (bool, *ShortURL, error) {
	originalShortURL := &ShortURL{}
	err := trans.SelectOne(originalShortURL, "SELECT * FROM urls WHERE original=?", shortURL.Original)
	if err == sql.ErrNoRows {
		return false, nil, nil
	} else if err != nil {
		return false, nil, err
	}

	return true, originalShortURL, nil
}

const APITokenSize = 20

type APIKey struct {
	Model
	CreatedAt int64  `db:"created_at"`
	Token     string `db:"token"`
}

func RandomToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (apiKey *APIKey) GenerateToken(trans *gorp.Transaction) error {
	var token string
	var err error

	for j := 0; j < 3; j++ {
		token, err = RandomToken(APITokenSize)
		if err != nil {
			return err
		}

		ret, err := trans.SelectInt("SELECT COUNT(*) FROM api_keys WHERE token=?", token)
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
