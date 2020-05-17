package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/astaxie/beego/logs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

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

// Retrieve remote tags without cloning repository
func main() {
	isEmpty := isRepoEmpty("git@vogsphere.msk.21-school.ru:vogsphere/intra-uuid-b82b46b0-0aab-4300-9bfb-d4f9765bfb74-2733882")
	fmt.Println(isEmpty)
}
