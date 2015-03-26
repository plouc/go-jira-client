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
    search_url = "/search/"
)


/*
    CreateIssue creates a JIRA issue with given IssueFields.
    Example:

    jira := gojira.NewJira(... args ...)

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
*/
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

    if fields.TimeTracking != nil {
        subData["timetracking"] = fields.TimeTracking
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
        fmt.Printf("%s\n%s\n",err,string(r))
	}

    return rsp
}

/**
  Search issues assigned to given user

  user          string
  maxResults    int
  startAt       int

  return IssueList
*/
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

/**
 Search issues using jira sql. Please see Jira documentation to know how to build queries

 jql            string  a JQL query string
 startAt        int     the index of the first issue to return (0-based)
 maxResults     int     the maximum number of issues to return (defaults to 50).
 validateQuery  bool    whether to validate the JQL query
 fields         string  the list of fields to return for each issue. By default, all navigable fields are returned.
 expand         string  A comma-separated list of the parameters to expand.

 return rsp     List of issues
*/
func (j *Jira) SearchIssues(jql string, startAt int,maxResults int,validateQuery bool, fields string, expand string) (rsp * IssueList){

    url := j.BaseUrl + j.ApiPath + issue_url + search_url

    url += "?jql=" + jql

    if startAt > 0 {
        url += fmt.Sprintf("&startAt=%d",startAt)
    }

    if maxResults > 0 {
        url += fmt.Sprintf("&maxResults=%d",maxResults)
    }

    if !validateQuery {
        url += fmt.Sprintf("&validateQuery=%t", validateQuery)
    }

    if fields != "" {
        url += "&fields=" + fields
    }

    if expand != "" {
        url += "&expand=" + expand
    }

    if j.Debug {
        fmt.Println(url)
    }

    result := j.buildAndExecRequest("GET", url, nil)

    err := json.Unmarshal(result, rsp)
    if err != nil {
        fmt.Println("%s", err)
    }

    if j.Debug {
        fmt.Println(result)
    }

    return rsp
}


/*
Search an issue by its id

id      string          Key id

return  Issue
*/
func (j *Jira) Issue(id string) (issue *Issue) {

    url := j.BaseUrl + j.ApiPath + issue_url + id
    contents := j.buildAndExecRequest("GET", url, nil)

    if j.Debug {
        fmt.Println(url)
    }

    err := json.Unmarshal(contents, &issue)
    if err != nil {
        fmt.Println("%s", err)
    }

    if j.Debug {
        fmt.Println(issue)
    }

    return issue
}

type IssueCreateResponse struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

type Issue struct {
    Id        string       `json:"id"`
    Key       string       `json:"key"`
	Self      string       `json:"self"`
	Expand    string       `json:"expand"`
	Fields    *IssueFields `json:"fields"`
	CreatedAt time.Time    `json:""`
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
	IssueType   *IssueType      `json:"issuetype"`
	Parent      *Issue          `json:""`
	Summary     string          `json:"summary"`
    Description string          `json:"description"`
	Reporter    *IssueUser      `json:"reporter"`
	Assignee    *IssueUser      `json:"assignee"`
	Project     *JiraProject    `json:""`
	Created     string          `json:"created"`
    TimeTracking *IssueTimeTracking `json:"timetracking"`
    Custom      map[string]interface{}
}

type IssueTimeTracking struct {
    OriginalEstimate    string `json:"originalEstimate"`
    RemainingEstimate    string `json:"remainingEstimate"`
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

