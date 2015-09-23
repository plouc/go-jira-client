package main

import (
	"os"
	"flag"
	"fmt"
	"time"
    "strings"
	//"github.com/plouc/go-jira-client"
	"github.com/jmervine/go-jira-client"
	"io/ioutil"
	"launchpad.net/goyaml"
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
        fmt.Println(err.Error())
        return
    }

    jira := gojira.NewJira(
        config.Host,
        config.ApiPath,
        config.ActivityPath,
        &gojira.Auth{config.Login, config.Password,},
    )

	var project string
	flag.StringVar(&project, "p", "", "[required] jira project")

	var summary string
	flag.StringVar(&summary, "s", "", "[required] issue summary")

	var desc string
	flag.StringVar(&desc, "d", "", "issue description")

	var assignee string
	flag.StringVar(&assignee, "a", "", "issue assignee username")

	var issuetype string
	flag.StringVar(&issuetype, "t", "Issue", "issue type")

    var custom string
    flag.StringVar(&custom, "c", "", "comma seperate key/value pairs of custom fields\n" +
                "\t> -c \"123451:1,123452:foo bar bah,123453:true\"")

	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *help == true || project == "" || summary == "" {
		flag.Usage()
		return
	}

    fields := &gojira.IssueFields{
        Project: &gojira.JiraProject{ Key: project },
        IssueType: &gojira.IssueType{ Name: issuetype },
        Summary: summary,
        Description: desc,
    }

    if assignee != "" {
		user, err := jira.User(assignee)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
        fields.Assignee = user.Assignee()
    }

    if custom != "" {
        fields.Custom = make(map[string]interface{})
        for _, i := range strings.Split(custom, ",") {
            s := strings.Split(i, ":")
            k := strings.TrimSpace(s[0])
            v := strings.TrimSpace(s[1])
            fields.Custom[k] = v
        }
    }

    if rsp := jira.CreateIssue(fields); rsp != nil {
        fmt.Printf("%+v\n", rsp)
    }
}
