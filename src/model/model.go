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
	Counter   int64  `db:"counter"`
}

func (shortURL *ShortURL) GenerateCode() error {
	if shortURL.Id != 0 {
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
	shortURL.Code = urlCode
	return nil
}

func FindShortURL(code string) (*ShortURL, error) {
	var shortURL ShortURL
	var errDeserialize error
	cacheKey := []byte("urls_" + code)

	bytes, errGet := database.Cache.Get(cacheKey)
	if errGet == nil {
		errDeserialize = util.Deserialize(bytes, &shortURL)
	}

	if bytes == nil || errGet != nil || errDeserialize != nil {
		err := database.DbMap.SelectOne(&shortURL, "SELECT * FROM urls WHERE code=?", code)
		if err != nil {
			return nil, err
		}

		bytes, err = util.Serialize(shortURL)
		if err == nil {
			_ = database.Cache.Set(cacheKey, bytes, 15*60)
		}
	}

	return &shortURL, nil
}

func SaveShortURL(shortURL *ShortURL) error {
	var err error

	if shortURL.Id == 0 {
		err = database.DbMap.Insert(shortURL)
	} else {
		_, err = database.DbMap.Update(shortURL)
	}

	if err != nil {
		return err
	}

	cacheKey := []byte("urls_" + shortURL.Code)
	bytes, err := util.Serialize(shortURL)
	if err != nil {
		_ = database.Cache.Set(cacheKey, bytes, 15*60)
	}

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

func FindAPIKey(token string) (*APIKey, error) {
	var apiKey APIKey
	var errDeserialize error
	cacheKey := []byte("api_keys_" + token)

	bytes, errGet := database.Cache.Get(cacheKey)
	if errGet == nil {
		errDeserialize = util.Deserialize(bytes, &apiKey)
	}

	if bytes == nil || errGet != nil || errDeserialize != nil {
		err := database.DbMap.SelectOne(&apiKey, "SELECT * FROM api_keys WHERE token=?", token)
		if err != nil {
			return nil, err
		}

		bytes, err = util.Serialize(apiKey)
		if err == nil {
			_ = database.Cache.Set(cacheKey, bytes, 5*60)
		}
	}

	return &apiKey, nil
}

func SaveAPIKey(apiKey *APIKey) error {
	var err error

	if apiKey.Id == 0 {
		err = database.DbMap.Insert(apiKey)
	} else {
		_, err = database.DbMap.Update(apiKey)
	}

	if err != nil {
		return err
	}

	cacheKey := []byte("api_keys_" + apiKey.Token)
	bytes, err := util.Serialize(apiKey)
	if err != nil {
		_ = database.Cache.Set(cacheKey, bytes, 5*60)
	}

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
