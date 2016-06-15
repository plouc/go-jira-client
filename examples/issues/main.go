package main

import (
	"flag"
	"fmt"
	"github.com/plouc/go-jira-client"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Host         string `yaml:"host"`
	ApiPath      string `yaml:"api_path"`
	ActivityPath string `yaml:"activity_path"`
	Login        string `yaml:"login"`
	Password     string `yaml:"password"`
}

func main() {
	startedAt := time.Now()
	defer func() {
		fmt.Printf("processed in %v\n", time.Now().Sub(startedAt))
	}()

	help := flag.Bool("help", false, "Show usage")

	// read config file
	file, e := ioutil.ReadFile("../config.yml")
	if e != nil {
		fmt.Printf("Config file error: %v\n", e)
		os.Exit(1)
	}

	// parse config file
	config := new(Config)
	err := goyaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	jira := gojira.NewJira(
		config.Host,
		config.ApiPath,
		config.ActivityPath,
		&gojira.Auth{config.Login, config.Password},
	)

	var method string
	flag.StringVar(&method, "m", "", "Specify method to retrieve issue(s) data, available methods:\n"+
		"  > -m issue -id ISSUE_ID")

	var id string
	flag.StringVar(&id, "id", "", "Specify issue id")

	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *help == true || method == "" {
		flag.Usage()
		return
	}

	switch method {
	case "issue":
		if id == "" {
			flag.Usage()
			return
		}

		issue := jira.Issue(id)

		format := "> %-14s: %s\n"

		fmt.Printf("%s\n", issue.Id)
		fmt.Printf(format, "self", issue.Self)
		fmt.Printf(format, "key", issue.Key)
		fmt.Printf(format, "expand", issue.Expand)
		fields := issue.Fields
		fmt.Printf(format, "summary", fields.Summary)
		fmt.Printf(format, "reporter", fields.Reporter.Name)
		fmt.Printf(format, "assignee", fields.Assignee.Name)
		fmt.Printf(format, "is subtask?", strconv.FormatBool(fields.IssueType.Subtask))
		//fmt.Printf(format, "created at", issue.CreatedAt)
	}
}
