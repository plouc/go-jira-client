package gojira_test

import (
    "fmt"
    "time"
    "testing"
    "encoding/json"
    . "github.com/jmervine/GoT"
    "github.com/hekima/go-jira-client"
)

func TestIssueCreateResponse(T *testing.T) {
    var raw_rsp = []byte(`{
        "id": "12345",
        "key": "ASDF",
        "self": "foobar"
    }`)

    var rsp gojira.IssueCreateResponse
    err := json.Unmarshal(raw_rsp, &rsp)

    Go(T).AssertNil(err)
    Go(T).AssertEqual(rsp.Id, "12345")
    Go(T).AssertEqual(rsp.Key, "ASDF")
    Go(T).AssertEqual(rsp.Self, "foobar")
}

func TestIssue(T *testing.T) {
    issue := gojira.Issue{
        Id: "ID",
        Key: "KEY",
        CreatedAt: time.Now(),
    }

    issue_json, err := json.Marshal(&issue)

    Go(T).AssertNil(err)
    Go(T).AssertEqual(issue_json, []byte(`{"key":"KEY"}`),
                "Issue should only marshal with Key")
}

func TestIssueType(T *testing.T) {
    it := gojira.IssueType{
        Self: "SELF",
        Id: "ID",
        Description: "DESCRIPTION",
        IconUrl: "ICON URL",
        Name: "NAME",
        Subtask: true,
    }

    it_json, err := json.Marshal(&it)

    Go(T).AssertNil(err)
    Go(T).AssertEqual(it_json, []byte(`{"name":"NAME","subtask":true}`),
                "Issue should only marshal with Name and Subtask")
}

func ExampleJira_CreateIssue() {
    jira := gojira.NewJira(
        "jira.host.com",
        "/rest/api/latest",
        "/activity",
        &gojira.Auth{"username", "password"},
    )

    // setting custom fields
    custom := make(map[string]interface{})
    custom["123451"] = 1
    custom["123452"] = "test custom data"

    fields := &gojira.IssueFields{
        Project: &gojira.JiraProject{ Key: "TEST" },
        Summary: "some new issue summary",
        Description: "some new issue description",
        IssueType: &gojira.IssueType{ Name: "bug" },
        Custom: custom,
    }

    if user, e := jira.User("someguy"); e == nil {
        fields.Assignee = user.Assignee()
    }

    rsp := jira.CreateIssue(fields)
    fmt.Printf("%+v\n", rsp)
}

