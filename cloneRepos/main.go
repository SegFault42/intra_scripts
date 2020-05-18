package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/astaxie/beego/logs"
	"github.com/go-git/go-git/v5"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type usersList []struct {
	RepoURL   string `json:"Repo URL"`
	Name      string `json:"Name"`
	FinalMark int    `json:"Final Mark"`
}

func getUsersList(jsonFile string) usersList {
	content, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var list usersList

	json.Unmarshal(content, &list)

	return list
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
		URL:  vogsphereRepo,
		Auth: sshAuth,
		// Progress: os.Stdout,
	})
	if err != nil {
		logs.Error(err)
	}
	os.RemoveAll(login + "/.git")
}

func parseArgs() *string {
	parser := argparse.NewParser("clone Repositories", "This program is made to clone many repository from 21vogsphere")
	args := parser.String("j", "json", &argparse.Options{Required: true, Help: "json to parse"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return nil
	}

	return args
}

func main() {
	args := parseArgs()
	if args == nil {
		return
	}

	list := getUsersList(*args)

	for i, _ := range list {
		list[i].Name = strings.Split(list[i].Name, "'")[0]
		if list[i].FinalMark > 0 {
			log.Println("Cloning ", list[i].Name, "'s repo")
			cloneRepo(list[i].RepoURL, list[i].Name)
		}
	}
}
