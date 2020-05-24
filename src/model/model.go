package model

import (
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"reflect"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/util"
)

var ErrAlreadyInDatabase = errors.New("already stored in database")
var ErrCannotBeGenerated = errors.New("cannot be generated")

func getModel(model interface{}, cacheKey string, cachePrefix string, cacheExpire int, query string, queryParams ...interface{}) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("getModel model must be a pointer, but got: %t", t)
	}

	cacheValue, err := database.Cache.Get([]byte(fmt.Sprintf("%s_%s", cachePrefix, cacheKey)))
	if err != nil && err != freecache.ErrNotFound {
		return err
	}

	if err == freecache.ErrNotFound {
		err = database.DbMap.SelectOne(model, query, queryParams...)
		if err != nil {
			return err
		}

		cacheValue, err = util.Serialize(model)
		if err != nil {
			return err
		}

		err = database.Cache.Set(cacheValue, cacheValue, cacheExpire)
		if err != nil {
			return err
		}
	} else {
		err = util.Deserialize(cacheValue, model)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveModel(model interface{}, modelId int64, cacheKey string, cachePrefix string, cacheExpire int) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("saveModel model must be a pointer, but got: %t", t)
	}

	var err error

	if modelId == 0 {
		err = database.DbMap.Insert(model)
	} else {
		_, err = database.DbMap.Update(model)
	}

	if err != nil {
		return err
	}

	cacheValue, err := util.Serialize(model)
	if err != nil {
		return err
	}

	return database.Cache.Set([]byte(fmt.Sprintf("%s_%s", cachePrefix, cacheKey)), cacheValue, cacheExpire)
}

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

func GetShortURL(code string) (*ShortURL, error) {
	shortURL := ShortURL{}
	err := getModel(
		&shortURL,
		code,
		"urls",
		15*60,
		"SELECT * FROM urls WHERE code=?",
		code,
	)
	if err != nil {
		return nil, err
	}

	return &shortURL, nil
}

func SaveShortURL(shortURL *ShortURL) error {
	return saveModel(
		shortURL,
		shortURL.Id,
		shortURL.Code,
		"urls",
		15*60,
	)
}

const APIKeySize = 20
const APIKeyChars = util.AsciiLetters + util.Digits + util.SpecialChars

type APIKey struct {
	Id    int64  `db:"id, primarykey, autoincrement"`
	Time  int64  `db:"time"`
	Token string `db:"token, size:20"`
}

func (apiKey *APIKey) GenerateToken() error {
	if apiKey.Id != 0 {
		return ErrAlreadyInDatabase
	}

	var token string

	for j := 0; j < 3; j++ {
		token = util.RandomString(APIKeySize, APIKeyChars)
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

func GetAPIKey(token string) (*APIKey, error) {
	apiKey := APIKey{}
	err := getModel(
		&apiKey,
		token,
		"api_keys",
		5*60,
		"SELECT * FROM api_keys WHERE token=?",
		token,
	)
	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func SaveAPIKey(apiKey *APIKey) error {
	return saveModel(
		apiKey,
		apiKey.Id,
		apiKey.Token,
		"api_keys",
		5*60,
	)
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
