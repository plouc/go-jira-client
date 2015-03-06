package gojira

import (
	"fmt"
	"time"
	"net/url"
	"encoding/xml"
	"encoding/json"
)

const (
	user_url        = "/user"
	user_search_url = "/user/search"
	// http://example.com:8080/jira/rest/api/2/user/assignable/multiProjectSearch [GET]
	// http://example.com:8080/jira/rest/api/2/user/assignable/search [GET]
	// http://example.com:8080/jira/rest/api/2/user/avatar [POST, PUT]
	// http://example.com:8080/jira/rest/api/2/user/avatar/temporary [POST, POST]
	// http://example.com:8080/jira/rest/api/2/user/avatar/{id} [DELETE]
	// http://example.com:8080/jira/rest/api/2/user/avatars [GET]
	// http://example.com:8080/jira/rest/api/2/user/picker [GET]
	// http://example.com:8080/jira/rest/api/2/user/viewissue/search [GET]
)

type User struct {
	Self         string            `json:"self"`
	Name         string            `json:"name"`
	EmailAddress string            `json:"emailAddress"`
	DisplayName  string            `json:"displayName"`
	Active       bool              `json:"active"`
	TimeZone     string            `json:"timeZone"`
	AvatarUrls   map[string]string `json:"avatarUrls"`
	Expand       string            `json:"expand"`
	// "groups": {
    //     "size": 3,
    //     "items": [
    //         {
    //             "name": "jira-user",
    //             "self": "http://www.example.com/jira/rest/api/2/group?groupname=jira-user"
    //         },
    //         {
    //             "name": "jira-admin",
    //             "self": "http://www.example.com/jira/rest/api/2/group?groupname=jira-admin"
    //         },
    //         {
    //             "name": "important",
    //             "self": "http://www.example.com/jira/rest/api/2/group?groupname=important"
    //         }
    //     ]
    // }
}

// Assignee is a helper method for abstracting IssueFields.Assignee
// when building data for CreateIssue
func (u *User) Assignee() *IssueUser {
    return &IssueUser{ Name: u.Name }
}

/*
Returns a user. This resource cannot be accessed anonymously.

    GET http://example.com:8080/jira/rest/api/2/user?username=USERNAME

Parameters

    username string The username

Usage

	user, err := jira.User("username")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%+v\n", user)
*/
func (j *Jira) User(username string) (*User, error) {
	url := j.BaseUrl + j.ApiPath + user_url + "?username=" + username
	contents := j.buildAndExecRequest("GET", url, nil)

	user := new(User)
	err := json.Unmarshal(contents, &user)
	if err != nil {
		fmt.Println("%s", err)
	}

	return user, err
}

/*
Returns a list of users that match the search string. This resource cannot be accessed anonymously.

	GET http://example.com:8080/jira/rest/api/2/user/search

Parameters

	username        string  A query string used to search username, name or e-mail address
	startAt         int     The index of the first user to return (0-based)
	maxResults      int     The maximum number of users to return (defaults to 50).
				   	        The maximum allowed value is 1000.
				   	        If you specify a value that is higher than this number,
				   	        your search results will be truncated.
	includeActive   boolean If true, then active users are included in the results (default true)
	includeInactive boolean If true, then inactive users are included in the results (default false)

*/
func (j *Jira) SearchUser(username string, startAt int, maxResults int, includeActive bool, includeInactive bool) {
	url := j.BaseUrl + j.ApiPath + user_url + "?username=" + username
	contents := j.buildAndExecRequest("GET", url, nil)
	fmt.Println(string(contents))

	// @todo
}


func (j *Jira) UserActivity(user string) (ActivityFeed, error) {
	url := j.BaseUrl + j.ActivityPath + "?streams=" + url.QueryEscape("user IS " + user)

	return j.Activity(url)
}

func (j *Jira) Activity(url string) (ActivityFeed, error) {

	contents := j.buildAndExecRequest("GET", url, nil)

	var activity ActivityFeed
	err := xml.Unmarshal(contents, &activity)
	if err != nil {
		fmt.Println("%s", err)
	}

	return activity, err
}


type ActivityItem struct {
	Title    string    `xml:"title"json:"title"`
	Id       string    `xml:"id"json:"id"`
	Link     []Link    `xml:"link"json:"link"`
	Updated  time.Time `xml:"updated"json:"updated"`
	Author   Person    `xml:"author"json:"author"`
	Summary  Text      `xml:"summary"json:"summary"`
	Category Category  `xml:"category"json:"category"`
}

type ActivityFeed struct {
	XMLName  xml.Name        `xml:"http://www.w3.org/2005/Atom feed"json:"xml_name"`
	Title    string          `xml:"title"json:"title"`
	Id       string          `xml:"id"json:"id"`
	Link     []Link          `xml:"link"json:"link"`
	Updated  time.Time       `xml:"updated,attr"json:"updated"`
	Author   Person          `xml:"author"json:"author"`
	Entries  []*ActivityItem `xml:"entry"json:"entries"`
}


type Category struct {
	Term string `xml:"term,attr"json:"term"`
}

type Link struct {
	Rel  string `xml:"rel,attr,omitempty"json:"rel"`
	Href string `xml:"href,attr"json:"href"`
}

type Person struct {
	Name     string `xml:"name"json:"name"`
	URI      string `xml:"uri"json:"uri"`
	Email    string `xml:"email"json:"email"`
	InnerXML string `xml:",innerxml"json:"inner_xml"`
}

type Text struct {
	Type string `xml:"type,attr,omitempty"json:"type"`
	Body string `xml:",chardata"json:"body"`
}
