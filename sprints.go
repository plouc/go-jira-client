package gojira

import (
	"fmt"
	"strconv"
	"encoding/json"
)

const (
	sprintQuery_url = "/sprintquery"
	sprintReport_url = "/rapid/charts/sprintreport"
)

type Sprint struct {
	Id        int           `json:"id"`
	Name      string        `json:"name"`
	StateKey  string        `json:"stateKey"`
	BoardName string        `json:"boardName"`
	State     string        `json:state`
	StartDate string        `json:startDate`
	EndDate   string        `json:endDate`
}

type Sprints struct {
	Sprints []Sprint `json:"sprints"`
}

type TextValue struct {
	Value float32        `json:"value"`
	Text  string        `json:"text"`
}

type SprintContents struct {
	CompletedIssuesEstimateSum TextValue        `json:"completedIssuesEstimateSum"`
	AllIssuesEstimateSum       TextValue        `json:"allIssuesEstimateSum"`
	ScopeChange                map[string]bool `json:"issueKeysAddedDuringSprint"`
}

type SprintReport struct {
	//	CompletedIssues	  []SprintIssues    `json:"completedIssues"`
	//	IncompletedIssues []SprintIssues	`json:"incompletedIssues"`
	Contents SprintContents    `json:"contents"`
	Sprint   Sprint            `json:"sprint"`

}

type SprintIssues struct {
	Id           int        `json:"id"`
	Key          string    `json:"key"`
	Hidden       bool    `json:"hidden"`
	TypeName     string    `json:"typeName"`
	Summary      string    `json:"sprint"`
	Assignee     string    `json:"sprint"`
	AssigneeName string    `json:"sprint"`
}


/*
 List sprints

 rapidViewId 	int		Id of board
 history 		bool 	List completed sprints
 future			bool	List sprints in the future

*/
func (j *Jira) ListSprints(rapidViewId int, history bool, future bool) (*Sprints, error) {

	url := j.BaseUrl + j.GreenHopper + sprintQuery_url + "/"+ strconv.Itoa(rapidViewId) + "?"

	url += "includeHistoricSprints="+strconv.FormatBool(history)
	url += "&includeFutureSprints="+strconv.FormatBool(future)

	if j.Debug {
		fmt.Println(url)
	}

	contents := j.buildAndExecRequest("GET", url, nil)

	sprints := new(Sprints)
	err := json.Unmarshal(contents, &sprints)
	if err != nil {
		fmt.Println("%s", err)
	}

	if j.Debug {
		fmt.Println(sprints)
	}

	return sprints, err
}

/*
  Get information of sprint report

  rapidViewId		int
  sprintId			int
 */
func (j *Jira) GetSprintReport(rapidViewId int, sprintId int) (*SprintReport, error) {

	url := j.BaseUrl + j.GreenHopper + sprintReport_url + "?"

	url += "rapidViewId="+strconv.Itoa(rapidViewId)
	url += "&sprintId="+strconv.Itoa(sprintId)

	if j.Debug {
		fmt.Println(url)
	}

	contents := j.buildAndExecRequest("GET", url, nil)

	sprintReport := new(SprintReport)
	err := json.Unmarshal(contents, &sprintReport)
	if err != nil {
		fmt.Println("%s", err)
	}

	if j.Debug {
		fmt.Println(sprintReport)
	}

	return sprintReport, err
}
