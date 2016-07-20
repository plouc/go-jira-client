package gojira

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
)

const (
	dateLayout = "2006-01-02T15:04:05.000-0700"
)

func NewJIRA(baseUrl string, auth *Auth) *Jira {
	client := &http.Client{}

	return &Jira{
		BaseUrl:      baseUrl,
		ApiPath:      "/rest/api/latest",
		ActivityPath: "/activity",
		GreenHopper:  "/rest/greenhopper/latest",
		Client:       client,
		Auth:         auth,
		Encoding:     "json",
	}
}

func NewJira(baseUrl string, apiPath string, activityPath string, auth *Auth) *Jira {

	client := &http.Client{}

	return &Jira{
		BaseUrl:      baseUrl,
		ApiPath:      apiPath,
		ActivityPath: activityPath,
		Client:       client,
		Auth:         auth,
	}
}

func (j *Jira) buildAndExecRequest(method string, url string, data io.Reader) []byte {

	req, err := http.NewRequest(method, url, data)
	if err != nil {
		panic("Error while building jira request")
	}
	req.SetBasicAuth(j.Auth.Login, j.Auth.Password)

	if data != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := j.Client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)

	return contents
}

func (p *Pagination) Compute() {
	p.PageCount = int(math.Ceil(float64(p.Total) / float64(p.MaxResults)))
	p.Page = int(math.Ceil(float64(p.StartAt) / float64(p.MaxResults)))

	p.Pages = make([]int, p.PageCount)
	for i := range p.Pages {
		p.Pages[i] = i
	}
}

type Jira struct {
	BaseUrl      string
	ApiPath      string
	ActivityPath string
	GreenHopper  string
	Encoding     string
	Debug        bool

	Auth   *Auth
	Client *http.Client
}

type Auth struct {
	Login    string
	Password string
}

type Pagination struct {
	Total      int
	StartAt    int
	MaxResults int
	Page       int
	PageCount  int
	Pages      []int
}

type JiraProject struct {
	Self       string            `json:"-"`
	Id         string            `json:"-"`
	Key        string            `json:"key"`
	Name       string            `json:"-"`
	AvatarUrls map[string]string `json:"-"`
}
