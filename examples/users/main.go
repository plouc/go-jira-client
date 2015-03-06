package main

import (
	"os"
	"flag"
	"fmt"
	"time"
	"strconv"
	"github.com/plouc/go-jira-client"
	"io/ioutil"
	"launchpad.net/goyaml"
	"github.com/hekima/go-jira-client"
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
        &gojira.Auth{config.Login, config.Password,},
    )

	var method string
	flag.StringVar(&method, "m", "", "Specify method to retrieve user(s) data, available methods:\n" +
									 "  > -m user -u USERNAME")

	var username string
	flag.StringVar(&username, "u", "", "Specify username")

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
	case "user":
		if username == "" {
			flag.Usage()
			return
		}

		user, err := jira.User(username)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		format := "> %-14s: %s\n"

		fmt.Printf("%s\n", user.Name)
		fmt.Printf(format, "self",          user.Self)
		fmt.Printf(format, "email address", user.EmailAddress)
		fmt.Printf(format, "display name",  user.DisplayName)
		fmt.Printf(format, "active",        strconv.FormatBool(user.Active))
		fmt.Printf(format, "time zone",     user.TimeZone)
	}
}
