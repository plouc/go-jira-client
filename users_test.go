package gojira_test

import (
    "fmt"
    "github.com/hekima/go-jira-client"
)

func ExampleJira_User() {
    jira := gojira.NewJira(
        "jira.host.com",
        "/rest/api/latest",
        "/activity",
        &gojira.Auth{"username", "password"},
    )
    user, err := jira.User("someuser")
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", user)
}
