package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/ihaiker/sudis/libs/errors"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/dao"
	"io/ioutil"
	"net/http"
	"time"
)

type GithubRelease struct {
	HTMLURL string `json:"html_url"`
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

var release *GithubRelease
var lastRefreshTime = time.Now()

func getRelease() (*GithubRelease, error) {
	if release != nil && time.Now().Before(lastRefreshTime.Add(time.Hour)) {
		return release, nil
	}
	lastRefreshTime = time.Now()
	release = &GithubRelease{HTMLURL: "", TagName: "", Name: ""}

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
			return release, err
		} else if response.StatusCode != 200 {
			return release, errors.New(response.Status)
		} else {
			defer response.Body.Close()
			if bs, err := ioutil.ReadAll(response.Body); err != nil {
				return release, err
			} else {
				return release, json.Unmarshal(bs, release)
			}
		}
	}
}

func dashboard(manger *cluster.DaemonManager) interface{} {
	return func() *dao.JSON {
		release, _ := getRelease()

		nodeProcess := manger.CacheAll()

		all := 0                    //所有陈谷
		started := 0                //已经启动的程序
		allNode := len(nodeProcess) //所有节点
		onlineNode := 0             //存活节点
		cpu := float64(0)           //总共使用CUP量
		ram := uint64(0)            //使用内存量

		for node, processes := range nodeProcess {
			if node.Status == dao.NodeStatusOnline {
				onlineNode += 1
			}
			all += len(processes)
			for _, process := range processes {
				if process.GetStatus().IsRunning() {
					cpu += process.Cpu
					ram += process.Rss
					started += 1
				}
			}
		}

		return &dao.JSON{
			"node": &dao.JSON{
				"all":    allNode,
				"online": onlineNode,
			},
			"process": &dao.JSON{
				"all":     all,
				"started": started,
			},
			"info": &dao.JSON{
				"CPU": fmt.Sprintf("%0.4f", cpu),
				"RAM": ram,
			},
			"version": release,
		}
	}
}
