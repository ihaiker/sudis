package http

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/ihaiker/sudis/master/dao"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
	"time"
)

type GithubRelease struct {
	HTMLURL string `json:"html_url"`
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

func getRelease() (*GithubRelease, error) {
	reslease := &GithubRelease{HTMLURL: "", TagName: "", Name: ""}
	if request, err := http.NewRequest("", "https://api.github.com/repos/ihaiker/sudis/releases/latest", nil); err != nil {
		return nil, err
	} else {
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Second * 5,
		}
		if response, err := client.Do(request); err != nil {
			return reslease, err
		} else if response.StatusCode != 200 {
			return reslease, errors.New(response.Status)
		} else {
			defer response.Body.Close()
			if bs, err := ioutil.ReadAll(response.Body); err != nil {
				return reslease, err
			} else {
				return reslease, json.Unmarshal(bs, reslease)
			}
		}
	}
}

func dashboard(ctx iris.Context) *JSON {

	nodes, err := dao.NodeDao.List()
	AssertErr(err)
	allNode := len(nodes)
	onlineNode := 0
	for _, node := range nodes {
		if node.Status == "online" {
			onlineNode += 1
		}
	}

	programs, err := dao.ProgramDao.List("", "", "", "", 1, 30000)
	AssertErr(err)
	allProgram := len(programs)
	startedProgram := 0
	for _, program := range programs {
		if program.Status.IsRunning() {
			startedProgram += 1
		}
	}
	release, _ := getRelease()

	return &JSON{
		"node": &JSON{
			"all":    allNode,
			"online": onlineNode,
		},
		"process": &JSON{
			"all":     allProgram,
			"started": startedProgram,
		},
		"info":    &JSON{"CPU": "Unrealized", "RAM": "Unrealized"},
		"version": release,
	}
}
