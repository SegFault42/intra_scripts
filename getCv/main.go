package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/go-git/go-git/v5"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type sToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	CreatedAt   int    `json:"created_at"`
}

type intraProject struct {
	Login         string `json:login`
	ProjectsUsers []struct {
		ID        int  `json:id`
		Validated bool `json:"validated?"`
		Project   struct {
			ID int `json:"id"`
		} `json:"project"`
	} `json:"projects_users"`
}

type vogsphereRepo struct {
	Teams []struct {
		RepoURL string `json:"repo_url"`
	} `json:"teams"`
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

func getProjectID(intraUid int, token string) (int, string) {
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
			return project.ProjectsUsers[i].ID, project.Login
		}
	}

	return 0, ""
}

func getVogsphereAddress(projectId int, token string) string {
	url := "https://api.intra.42.fr/v2/projects_users/" + strconv.Itoa(projectId)
	resp := bearerGetRequest(token, url)

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var vogsphere vogsphereRepo

	err := json.Unmarshal([]byte(body), &vogsphere)
	if err != nil {
		panic(err)
	}

	return vogsphere.Teams[0].RepoURL
}

func cloneRepo(vogsphereRepo, login string) {
	dir := "./" + login
	os.Mkdir(dir, os.ModePerm)

	// clone with ssh key
	currentUser, err := user.Current()
	if err != nil {
		logs.Error(err)
	}
	sshAuth, err := ssh.NewPublicKeysFromFile("git", currentUser.HomeDir+"/.ssh/id_rsa", "admin")
	if err != nil {
		logs.Error(err)
	}
	// Clones the repository into the given dir, just as a normal git clone does
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:      vogsphereRepo,
		Auth:     sshAuth,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatal(err)
	}
	os.RemoveAll(login + "/.git")

	// Prints the content of the CHANGELOG file from the cloned repository
	// changelog, err := os.Open(filepath.Join(dir, ".git/config"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// io.Copy(os.Stdout, changelog)
}

func main() {
	clientId := os.Getenv("CLIENTID")
	clientSecret := os.Getenv("CLIENTSECRET")

	token := getNewToken(clientId, clientSecret)
	log.Println("Get token OK")

	projectId, login := getProjectID(61663, token)
	log.Println("GetProjectId ok")

	if projectId != 0 {
		vogsphereAddress := getVogsphereAddress(projectId, token)
		log.Println("Get Vogsphere address ok")
		if vogsphereAddress != "" {
			log.Println("Cloning Repo ...")
			cloneRepo(vogsphereAddress, login)
		}
	}
}
