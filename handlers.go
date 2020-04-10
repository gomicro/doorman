package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gomicro/doorman/users"
)

type token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func handleUserInfo(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "jwt ")
	log.Debugf("Token: %v", token)

	u, err := users.Lookup(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request token"))
		return
	}

	b, err := json.Marshal(u)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal user: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

func handleGetGoogleAuth(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	cID := q.Get("client_id")
	log.Debugf("Client ID: %v", cID)

	rType := q.Get("response_type")
	log.Debugf("Response Type: %v", rType)

	scope := q.Get("scope")
	log.Debugf("Scope: %v", scope)

	encState := q.Get("state")
	state, err := base64.StdEncoding.DecodeString(encState)
	if err != nil {
		log.Errorf("failed to decode state: %v", err.Error())

		w.WriteHeader(http.StatusBadRequest)

		msg := "state not base64 encoded"
		w.Write([]byte(msg))
		return
	}

	log.Debugf("State: %v", string(state))

	redURI := q.Get("redirect_uri")
	log.Debugf("Redirect URI: %v", redURI)

	ru, err := url.Parse(redURI)
	if err != nil {
		log.Errorf("failed to parse redirect uri: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad redirect uri"))
		return
	}

	code := "somecode"

	values := url.Values{}
	values.Set("code", code)
	values.Set("state", encState)

	ru.RawQuery = values.Encode()

	http.Redirect(w, r, ru.String(), http.StatusSeeOther)
}

func handlePostGoogleAuth(w http.ResponseWriter, r *http.Request) {
	u := users.Random()

	t := &token{
		AccessToken:  u.Sub,
		TokenType:    "jwt",
		RefreshToken: "somerefreshtoken",
		ExpiresIn:    time.Now().Add(15 * time.Minute).Unix(),
	}

	b, err := json.Marshal(t)
	if err != nil {
		log.Errorf("failed to marshal token: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to marshal token"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
