package gojira

import (
	"encoding/json"
	"fmt"
	"time"
)

// Do a get request in this url
//
// Timesheet api doc http://www.jiratimesheet.com/wiki/RESTful_endpoint.html
// https://zahpee.atlassian.net/rest/timesheet-gadget/1.0/raw-timesheet.json?targetUser=gabriel.campos&startDate=2015-08-01

const (
	worklog_url       = "/rest/api/latest/search"
	worklogDateFormat = "2006-01-02"
)

type WorklogResponse struct {
	Worklog []Worklog `json:"worklog"`
}

type Worklog struct {
	IssueKey     string `json:"key"`
	IssueSummary string `json:"summary"`
	Entries      []Logs `json:"entries"`
}

type Logs struct {
	Id               int    `json:"id"`
	Comment          string `json:"comment"`
	TimeSpent        int    `json:"timeSpent"`
	Author           string `json:"author"`
	AuthorName       string `json:"authorFullName"`
	Created          int    `json:"created"`
	StartDate        int    `json:"startDate"`
	UpdateAuthor     string `json:"updateAuthor"`
	UpdateAuthorName string `json:"updateAuthorFullName"`
	Updated          int    `json:"updated"`
}

func (j *Jira) Worklog(username string, start, end time.Time) (*WorklogResponse, error) {
	url := j.BaseUrl + worklog_url + "?targetUser=" + username + "&worklogDate=" + start.Format(worklogDateFormat) + "&worklogDate=" + end.Format(worklogDateFormat)
	url += "fields=worklog"

	if j.Debug {
		fmt.Println(url)
	}

	contents := j.buildAndExecRequest("GET", url, nil)

	response := new(WorklogResponse)
	er := json.Unmarshal(contents, &response)
	if er != nil {
		fmt.Println("%s", er)
	}

	if j.Debug {
		fmt.Println(response)
	}

	return response, nil
}
