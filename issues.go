package gojira

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	issue_worklog_url = "/worklog"
	issue_url         = "/issue"
	search_url        = "/search"
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
func (j *Jira) CreateIssue(fields *IssueFields) (rsp IssueCreateResponse) {

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
	err = json.Unmarshal(r, &rsp)
	if err != nil {
		fmt.Printf("%s\n%s\n", err, string(r))
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
	//
	//    for _, issue := range issues.Issues {
	//        t, _ := time.Parse(dateLayout, issue.Fields.Created)
	//        issue.CreatedAt = t
	//    }
	//
	//    pagination := Pagination{
	//        Total:      issues.Total,
	//        StartAt:    issues.StartAt,
	//        MaxResults: issues.MaxResults,
	//    }
	//    pagination.Compute()
	//
	//    issues.Pagination = &pagination

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
func (j *Jira) SearchIssues(jql string, startAt int, maxResults int, validateQuery bool, fields string, expand string) (rsp IssueList) {

	requestUrl := j.BaseUrl + j.ApiPath + search_url

	requestUrl += "?jql=" + url.QueryEscape(jql)

	if startAt > 0 {
		requestUrl += fmt.Sprintf("&startAt=%d", startAt)
	}

	if maxResults > 0 {
		requestUrl += fmt.Sprintf("&maxResults=%d", maxResults)
	}

	if !validateQuery {
		requestUrl += fmt.Sprintf("&validateQuery=%t", validateQuery)
	}

	if fields != "" {
		requestUrl += "&fields=" + fields
	}

	if expand != "" {
		requestUrl += "&expand=" + expand
	}

	if j.Debug {
		fmt.Println(requestUrl)
	}

	result := j.buildAndExecRequest("GET", requestUrl, nil)
	//print(string(result[:]))
	err := json.Unmarshal(result, &rsp)
	if err != nil {
		fmt.Println("ERR: %s", err)
	}

	if j.Debug {
		//        fmt.Println(result)
		//        fmt.Println(rsp)
	}

	return rsp
}

// Adds a new worklog entry to an issue.
//
// https://jira.atlassian.com/plugins/servlet/restbrowser#/resource/api-2-issue-issueidorkey-worklog/POST
//
// issue - the worklogs belongs to
// adjust - (optional) allows you to provide specific instructions to update the remaining time estimate of the issue.
// 			Valid values are
// 			"new" - sets the estimate to a specific value
// 			"leave"- leaves the estimate as is
//			"manual" - specify a specific amount to increase remaining estimate by
// 			"auto"- Default option. Will automatically adjust the value based on the new timeSpent specified on the
// 			worklog
// new - 	(required when "new" is selected for adjustEstimate) the new value for the remaining estimate field. "2d"
// reduceBy - (required when "manual" is selected for adjustEstimate) the amount to reduce the remaining estimate by "2d"
func (j *Jira) LogWork(issue, adjust, new, reduceBy string, worklog IssueWorklog) {

	requestUrl := j.BaseUrl + j.ApiPath + issue_url + "/" + issue + issue_worklog_url + "?adjustEstimate=auto"

	fmt.Println(requestUrl)

	requestBody, err := json.Marshal(worklog)
	if err != nil {
		fmt.Println("ERR: %s", err)
	}

	fmt.Println(string(requestBody))

	contents := j.buildAndExecRequest("POST", requestUrl, strings.NewReader(string(requestBody)))

	fmt.Println(string(contents[:]))
}

/*
Search an issue by its id

id      string          Key id

return  Issue
*/
func (j *Jira) Issue(id string) (issue *Issue) {

	url := j.BaseUrl + j.ApiPath + issue_url + "/" + id
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

func (j *Jira) IssuesByRawJQL(jql string) IssueList {
	url := j.BaseUrl + j.ApiPath + "/search?jql=" + url.QueryEscape(jql)
	return j.queryToIssueList(url)
}

func (j *Jira) queryToIssueList(url string) IssueList {
	contents := j.buildAndExecRequest("GET", url, nil)

	var issues IssueList
	err := json.Unmarshal(contents, &issues)
	if err != nil {
		fmt.Println("%s", err)
	}

	pagination := Pagination{
		Total:      issues.Total,
		StartAt:    issues.StartAt,
		MaxResults: issues.MaxResults,
	}
	pagination.Compute()

	return issues
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
	CreatedAt time.Time
	ChangeLog *IssueChangeLog `json:"changelog"`
}

type IssueChangeLog struct {
	StartAt   int                      `json:"startAt"`
	MaxResult int                      `json:"maxResults"`
	Total     int                      `json:"total"`
	Histories []*IssueChangeLogHistory `json:"histories"`
}

type IssueChangeLogHistory struct {
	Id      string                   `json:"id"`
	Author  *IssueAuthor             `json:"author"`
	Created CustomTime               `json:"created"`
	Item    []map[string]interface{} `json:"items"`
}

type IssueList struct {
	Expand     string      `json:"expand"`
	StartAt    int         `json:"startAt"`
	MaxResults int         `json:"maxResults"`
	Total      int         `json:"total"`
	Issues     []*Issue    `json:"issues"`
	Pagination *Pagination `json:"id"`
}

type IssueStatus struct {
	Self        string `json:"self"`
	Description string `json:"description"`
	Icon        string `json:"iconUrl"`
	Name        string `json:"name"`
}

type IssueUser struct {
	Name string `json:"name"`
}

type IssuePriority struct {
	Id   string `json:id`
	Self string `json:self`
	//    IconUrl string      `json:id`
	Name string `json:name`
}

type WorkLogList struct {
	StartAt    int        `json:"startAt"`
	MaxResults int        `json:"maxResults"`
	Total      int        `json:"total"`
	WorkLogs   []WorkLogs `json:"worklogs"`
}

type WorkLogs struct {
	Id               string       `json:"id"`
	Self             string       `json:"self"`
	Comment          string       `json:"comment"`
	TimeSpent        string       `json:"timeSpent"`
	TimeSpentSeconds int          `json:"timeSpentSeconds"`
	Author           *IssueAuthor `json:"author"`
	//AuthorName       string     `json:"authorFullName"`
	//Created          int        `json:"created"`
	StartDate CustomTime `json:"started"`
	//UpdateAuthor     string     `json:"updateAuthor"`
	//UpdateAuthorName string     `json:"updateAuthorFullName"`
	//Updated          int        `json:"updated"`
}

type IssueFields struct {
	IssueType    *IssueType         `json:"issuetype"`
	Parent       *Issue             `json:"parent"`
	Summary      string             `json:"summary"`
	Description  string             `json:"description"`
	Reporter     *IssueUser         `json:"reporter"`
	Assignee     *IssueUser         `json:"assignee"`
	Project      *JiraProject       `json:"project"`
	Priority     *IssuePriority     `json:"priority"`
	Created      CustomTime         `json:"created"`
	TimeSpent    int                `json:"timespent"`
	TimeEstimate int                `json:"aggregatetimeoriginalestimate"`
	TimeTracking *IssueTimeTracking `json:"timetracking"`
	Status       *IssueStatus       `json:"status"`
	SprintPoints float32            `json:"customfield_10004"`
	Labels       []string           `json:"labels"`
	WorkLog      *WorkLogList       `json:"worklog"`
	Custom       map[string]interface{}
}

type IssueTimeTracking struct {
	OriginalEstimate  string `json:"originalEstimate"`
	RemainingEstimate string `json:"remainingEstimate"`
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

type IssueAuthor struct {
	Self         string `json:"self"`
	Name         string `json:"name"`
	Key          string `json:"key"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
}

type IssueWorklog struct {
	Self          string      `json:"self"`
	Comment       string      `json:"comment"`
	TimeSpent     string      `json:"timeSpent"`
	Author        IssueAuthor `json:"author"`
	UpdatedAuthor IssueAuthor `json:"updateAuthor"`
}
