package storage

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/fatihdumanli/cnote/pkg/oauthv2"
	"github.com/xujiajun/nutsdb"
)

var InvalidTokenType = errors.New("Token type is invalid")

const TOKEN_KEY = "msgraphtoken"
const BUCKET = "cnote"

type TokenStatus int

const (
	DoesntExist TokenStatus = iota
	Expired
	Valid
)

//expired, doesnt exist
func CheckToken() (oauthv2.OAuthToken, TokenStatus) {

	var token oauthv2.OAuthToken

	db, closer, err := openDb()
	defer closer()

	//TODO: perhaps a better handling is required.
	if err != nil {
		return token, DoesntExist
	}
	var e *nutsdb.Entry

	fnGet := func(tx *nutsdb.Tx) error {
		key := []byte(TOKEN_KEY)

		if e, err = tx.Get(BUCKET, key); err != nil {
			return err
		}

		return nil
	}
	err = db.View(fnGet)

	//TODO: perhaps a better handling is required.
	if err != nil {
		return token, DoesntExist
	}

	err = json.Unmarshal(e.Value, &token)

	//TODO: perhaps a better handling is required.
	if err != nil {
		log.Fatal(err)
		return token, DoesntExist
	}

	//check if it expired
	if token.IsExpired() {
		return token, Expired
	} else {
		return token, Valid
	}
}

func StoreToken(t interface{}) error {

	if _, ok := t.(oauthv2.OAuthToken); !ok {
		return InvalidTokenType
	}

	//open nuts db
	db, closer, err := openDb()
	defer closer()

	if err != nil {
		log.Fatal(err)
		return err
	}

	//convert the token into bytes
	bytes, err := json.Marshal(t)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fnUpdate := func(tx *nutsdb.Tx) error {
		key := []byte(TOKEN_KEY)
		val := bytes
		if err := tx.Put(BUCKET, key, val, 0); err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	}

	//save the token
	err = db.Update(fnUpdate)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

//opens the nuts db
//returns nuts db, closer and an error
//call closer to clean up the resources.
func openDb() (*nutsdb.DB, func() error, error) {
	opts := nutsdb.DefaultOptions
	opts.Dir = "/tmp/cnotedb"
	db, err := nutsdb.Open(opts)

	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	return db, db.Close, nil
}
