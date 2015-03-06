package gojira

import (
	"fmt"
	"strconv"
	"encoding/json"
)

const (
	sprintQuery_url		= "/sprintquery"
	sprintReport_url	= "/rapid/charts/sprintreport"
)

type Sprint struct {
	Id		 	int			`json:"id"`
	Name	 	string		`json:"name"`
	StateKey	string		`json:"stateKey"`
	BoardName 	string		`json:"boardName"`
	State		string		`json:state`
}

type Sprints struct {
	Sprints []Sprint `json:"sprints"`
}

type TextValue struct {
	Value 		float32		`json:"value"`
	Text		string		`json:"text"`
}

type SprintContents struct {
	CompletedIssuesEstimateSum 	TextValue	`json:"completedIssuesEstimateSum"`
	AllIssuesEstimateSum		TextValue	`json:"allIssuesEstimateSum"`
}

type SprintReport struct {
	Contents	SprintContents	`json:"contents"`
	Sprint		Sprint			`json:"sprint"`
}

//allIssuesEstimateSum

/*
 List sprints

 rapidViewId 	int		Id of board
 history 		bool 	List completed sprints
 future			bool	List sprints in the future

*/
func (j *Jira) ListSprints(rapidViewId int,history bool,future bool) (*Sprints, error) {

	url := j.BaseUrl + j.GreenHopper + sprintQuery_url + "/"+ strconv.Itoa(rapidViewId) + "?"

	url += "includeHistoricSprints="+strconv.FormatBool(history)
	url += "&includeFutureSprints="+strconv.FormatBool(future)

//	fmt.Println(url)

	contents := j.buildAndExecRequest("GET", url)

	sprints := new(Sprints)
	err := json.Unmarshal(contents, &sprints)
	if err != nil {
		fmt.Println("%s", err)
	}

	return sprints, err
}

/*
  Get information of sprint report

  rapidViewId		int
  sprintId			int
 */
func (j *Jira) GetSprintReport(rapidViewId int,sprintId int) (*SprintReport, error) {

	url := j.BaseUrl + j.GreenHopper + sprintReport_url + "?"

	url += "rapidViewId="+strconv.Itoa(rapidViewId)
	url += "&sprintId="+strconv.Itoa(sprintId)

//	fmt.Println(url)

	contents := j.buildAndExecRequest("GET", url)

	sprintReport := new(SprintReport)
	err := json.Unmarshal(contents, &sprintReport)
	if err != nil {
		fmt.Println("%s", err)
	}

	return sprintReport, err
}
