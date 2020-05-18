package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/user"

	"github.com/astaxie/beego/logs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type intraTeams []struct {
	RepoURL   string `json:"Repo URL"`
	Login     string `json:"Name"`
	ID        int    `json:"ID"`
	ProjectID int    `json:"Project ID"`
}

func setSshAuth() (*ssh.PublicKeys, error) {
	// clone with ssh key
	currentUser, err := user.Current()
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	sshAuth, err := ssh.NewPublicKeysFromFile("git", currentUser.HomeDir+"/.ssh/id_rsa", "admin")
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return sshAuth, err
}

func isRepoEmpty(repo string) bool {
	// Create the remote with repository URL
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{repo},
	})

	// setup ssh Authentification
	sshAuth, err := setSshAuth()
	if err != nil {
		log.Panic(err)
	}

	// We can then use every Remote functions to retrieve wanted information
	_, err = rem.List(&git.ListOptions{
		Auth: sshAuth,
	})
	if err != nil {
		return true
		logs.Error(err)
	}

	return false
}

func parseJson(jsonFile string) (intraTeams, error) {
	var list intraTeams

	content, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return list, err
	}

	json.Unmarshal(content, &list)

	return list, nil
}

// Retrieve remote tags without cloning repository
func main() {
	vogsphereList, err := parseJson("IntraTeams.json")
	if err != nil {
		log.Fatal(err)
	}
	for i, _ := range vogsphereList {
		isEmpty := isRepoEmpty(vogsphereList[i].RepoURL)
		logs.Info("\nLogin: ", vogsphereList[i].Login,
			"\nRepoURL: ", vogsphereList[i].RepoURL,
			"\nProjectID: ", vogsphereList[i].ProjectID,
			"\nID: ", vogsphereList[i].ID,
			"\nIs empty: ", isEmpty)
	}

}
