package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type sToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	CreatedAt   int    `json:"created_at"`
}

type intraProject struct {
	ProjectsUsers []struct {
		ID        int  `json:id`
		Validated bool `json:"validated?"`
		Project   struct {
			ID int `json:"id"`
		} `json:"project"`
	} `json:"projects_users"`
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

func bearerGetRequest(token, url string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func getProjectID(intraUid int, token string) int {
	url := "https://api.intra.42.fr/v2/users/" + strconv.Itoa(intraUid)
	resp := bearerGetRequest(token, url)

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var project intraProject

	err := json.Unmarshal([]byte(body), &project)
	if err != nil {
		panic(err)
	}

	// check if project 902 is validate
	for i := 0; i < len(project.ProjectsUsers); i++ {
		if project.ProjectsUsers[i].Project.ID == 902 && project.ProjectsUsers[i].Validated == true {
			return project.ProjectsUsers[i].ID
		}
	}

	return 0
}

func main() {
	clientId := os.Getenv("CLIENTID")
	clientSecret := os.Getenv("CLIENTSECRET")

	token := getNewToken(clientId, clientSecret)

	projectId := getProjectID(61663, token)
	fmt.Println(projectId)
}
