package gojira_test

import (
	"encoding/json"
	. "github.com/franela/goblin"
	"github.com/hekima/go-jira-client"
	"testing"
)

func TestIssueCreateResponse(T *testing.T) {
	var raw_rsp = []byte(`{
        "id": "12345",
        "key": "ASDF",
        "self": "foobar"
    }`)

	var rsp gojira.IssueCreateResponse
	err := json.Unmarshal(raw_rsp, &rsp)

	g := Goblin(T)
	g.Assert(err).Equal(nil)
	g.Assert(rsp.Id).Equal("12345")
	g.Assert(rsp.Key).Equal("ASDF")
	g.Assert(rsp.Self).Equal("foobar")
}

func TestIssue(T *testing.T) {
	issue := gojira.Issue{
		Id:  "ID",
		Key: "KEY",
	}

	issue_json, err := json.Marshal(&issue)

	g := Goblin(T)
	g.Assert(err).Equal(nil)
	g.Assert(issue_json).Equal([]byte(`{"id":"ID","key":"KEY","self":"","expand":"","fields":null,"CreatedAt":"0001-01-01T00:00:00Z","changelog":null}`))
}

func TestIssueType(T *testing.T) {
	it := gojira.IssueType{
		Self:        "SELF",
		Id:          "ID",
		Description: "DESCRIPTION",
		IconUrl:     "ICON URL",
		Name:        "NAME",
		Subtask:     true,
	}

	it_json, err := json.Marshal(&it)

	g := Goblin(T)
	g.Assert(err).Equal(nil)
	g.Assert(it_json).Equal([]byte(`{"name":"NAME","subtask":true}`))
}
