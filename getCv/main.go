package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type info struct {
	id    string `json:"id"`
	login string `json:"login"`
}

type sToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	CreatedAt   int    `json:"created_at"`
}

func getNewToken(uid, secret string) string {
	reader := strings.NewReader(`grant_type=client_credentials&client_id=` + uid + `&client_secret=` + secret)
	req, err := http.NewRequest("POST", "https://api.intra.42.fr/oauth/token", reader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var tokenJson sToken
	err = json.Unmarshal([]byte(body), &tokenJson)
	if err != nil {
		panic(err)
	}

	return (tokenJson.AccessToken)
}

func intraUidToLogin(uid int) {
	url := "https://api.intra.42.fr/v2/users/45317"

	// get intra id
	var clientId string = os.Getenv("CLIENTID")
	var clientSecret string = os.Getenv("CLIENTSECRET")

	token := getNewToken(clientId, clientSecret)

	// get user list project
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))

	// var user info
	// err = json.Unmarshal([]byte(string(body)), &user)
	// if err != nil {
	// 	panic(err)
	// 	return
	// }
	// fmt.Println(user)
}

func main() {
	intraUidToLogin(50721)
}