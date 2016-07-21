/*
Package gojira is a work in progress lib for accessing JIRA's via REST.

    jira := gojira.NewJira(
        "jira.host.com",
        "/rest/api/latest",
        "/activity",
        &gojira.Auth{"username", "password"},
    )

TODO:
* More tests.
* GET Issue(s) information.
* Update Issue

*/
package gojira
