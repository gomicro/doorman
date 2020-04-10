package users

import (
	"fmt"
	"math/rand"
	"strconv"
)

type GoogleUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

var (
	store map[string]*GoogleUser
)

func init() {
	store = make(map[string]*GoogleUser)
}

func Lookup(token string) (*GoogleUser, error) {
	u, found := store[token]
	if !found {
		return nil, fmt.Errorf("failed to find user")
	}

	return u, nil
}

func Random() *GoogleUser {
	i := 1000000000000000 + rand.Int63()
	sub := strconv.FormatInt(i, 10)
	profile := fmt.Sprintf("https://plus.google.com/%v", sub)
	pic := fmt.Sprintf("https://lh5.googleusercontent.com/%v/photo.jpg", sub)

	u := &GoogleUser{
		Sub:           sub,
		Name:          "Daniel Hess",
		GivenName:     "Daniel",
		FamilyName:    "Hess",
		Profile:       profile,
		Picture:       pic,
		Email:         "dan9186@gmail.com",
		EmailVerified: true,
	}

	store[sub] = u
	return u
}
