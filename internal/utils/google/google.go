package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"kedai/backend/be-kedai/config"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}

func ValidateGoogleToken(tokenString string) (model.GoogleClaim, error) {
	claimsStruct := model.GoogleClaim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}
			return key, nil
		},
	)
	if err != nil {
		return model.GoogleClaim{}, err
	}

	claims, ok := token.Claims.(*model.GoogleClaim)
	if !ok {
		return model.GoogleClaim{}, errors.New("invalid Google JWT")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return model.GoogleClaim{}, errors.New("iss is invalid")
	}

	if claims.Audience[0] != config.GetEnv("GOOGLE_CLIENT_ID", "") {
		return model.GoogleClaim{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt.Unix() < time.Now().UTC().Unix() {
		return model.GoogleClaim{}, errors.New("JWT is expired")
	}

	return *claims, nil
}
