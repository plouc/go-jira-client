package gojira

import (
    "fmt"
    "strings"
    "net/url"
	"encoding/json"
    "strconv"
    "time"
)

const (
    issue_url = "/issue/"
)


// CreateIssue creates a JIRA issue with given IssueFields.
// Example:
//
//      jira := gojira.NewJira(... args ...)
//
//      // setting custom fields
//      custom := make(map[string]interface{})
//      custom["123451"] = 1
//      custom["123452"] = "test custom data"
//
//      fields := &gojira.IssueFields{
//          Project: &gojira.JiraProject{ Key: "TEST" },
//          Summary: "some new issue summary",
//          Description: "some new issue description",
//          IssueType: &gojira.IssueType{ Name: "bug" },
//          Custom: custom,
//      }
//
//      if user, e := jira.User("someguy"); e == nil {
//          fields.Assignee = user.Assignee()
//      }
//
//      rsp := jira.CreateIssue(fields)
//
func (j *Jira) CreateIssue(fields *IssueFields) (rsp *IssueCreateResponse) {

    // Support custom fields.
    dynData := make(map[string]interface{})
    subData := make(map[string]interface{})

    // Required
    if fields.IssueType != nil && fields.IssueType.Name != "" {
        subData["issuetype"] = fields.IssueType
    } else {
        fmt.Println("Error CreateIssue requires *IssueFields.IssueType.Name")
        return
    }

    // Required
    if fields.Project != nil && fields.Project.Key != "" {
        subData["project"] = fields.Project
    } else {
        fmt.Println("Error CreateIssue requires *IssueFields.Project.Key")
        return
    }

    if fields.Parent != nil {
        subData["parent"] = fields.Parent
    }

    if fields.Assignee != nil {
        subData["assignee"] = fields.Assignee
    }

	for k, v := range fields.Custom {
		subData["custom_"+k] = v
	}

    if fields.Summary == "" {
        fmt.Println("Error CreateIssue requires *IssueFields.Summary")
        return
    }
    subData["summary"] = fields.Summary

    if fields.Description != "" {
        subData["description"] = fields.Description
    }

    var postData []byte
    var err error

    url := j.BaseUrl + j.ApiPath + issue_url

    dynData["fields"] = subData
    postData, err = json.Marshal(dynData)
	if err != nil {
        fmt.Printf("Error marshaling fields: %s\n", err)
        return
	}

    r := j.buildAndExecRequest("POST", url, strings.NewReader(string(postData)))
    err = json.Unmarshal(r, rsp)
	if err != nil {
        fmt.Printf("Error unmarshaling response: %s\n\n", err)
        fmt.Println("Raw request:")
        fmt.Printf("%q\n\n", string(postData))
        fmt.Println("Raw response:")
        fmt.Printf("%s\n\n", string(r))
	}

    return
}

// search issues assigned to given user
func (j *Jira) IssuesAssignedTo(user string, maxResults int, startAt int) IssueList {

    url := j.BaseUrl + j.ApiPath + "/search?jql=assignee=\"" + url.QueryEscape(user) + "\"&startAt=" + strconv.Itoa(startAt) + "&maxResults=" + strconv.Itoa(maxResults)
    contents := j.buildAndExecRequest("GET", url, nil)

    var issues IssueList
    err := json.Unmarshal(contents, &issues)
    if err != nil {
        fmt.Println("%s", err)
    }

    for _, issue := range issues.Issues {
        t, _ := time.Parse(dateLayout, issue.Fields.Created)
        issue.CreatedAt = t
    }

    pagination := Pagination{
        Total:      issues.Total,
        StartAt:    issues.StartAt,
        MaxResults: issues.MaxResults,
    }
    pagination.Compute()

    issues.Pagination = &pagination

    return issues
}

// search an issue by its id
func (j *Jira) Issue(id string) Issue {

    url := j.BaseUrl + j.ApiPath + "/issue/" + id
    contents := j.buildAndExecRequest("GET", url, nil)

    var issue Issue
    err := json.Unmarshal(contents, &issue)
    if err != nil {
        fmt.Println("%s", err)
    }

    return issue
}

type IssueCreateResponse struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

type Issue struct {
    Id        string       `json:"-"`
    Key       string       `json:"key"`
	Self      string       `json:"-"`
	Expand    string       `json:"-"`
	Fields    *IssueFields `json:"-"`
	CreatedAt time.Time    `json:"-"`
}

type IssueList struct {
	Expand     string
	StartAt    int
	MaxResults int
	Total      int
	Issues     []*Issue
	Pagination *Pagination
}

type IssueUser struct {
	Name string `json:"name"`
}

type IssueFields struct {
	IssueType   *IssueType
	Parent      *Issue
	Summary     string
	Description string
	Reporter    *IssueUser
	Assignee    *IssueUser
	Project     *JiraProject
	Created     string
    Custom      map[string]interface{}
}

// IssueType is mainly used for creating an issue, as such we're only
// including json mapping for that which we want added to the out going
// json message.
type IssueType struct {
    Self        string `json:"-"`
    Id          string `json:"-"`
    Description string `json:"-"`
    IconUrl     string `json:"-"`
    Name        string `json:"name"`
    Subtask     bool   `json:"subtask"`
}

