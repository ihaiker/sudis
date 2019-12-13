package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/conf"
	"github.com/ihaiker/sudis/daemon"
	"strconv"
	"strings"
	"time"
)

func MakeServerCommand(dm *daemon.DaemonManager) rpc.OnMessage {
	var OK = []byte("OK")

	return func(channel remoting.Channel, request *rpc.Request) *rpc.Response {
		logger.Debug("action: ", request.URL, ", args:", string(request.Body))
		response := rpc.NewResponse(request.ID())

		switch request.URL {
		case "start":
			args := []string{}
			if response.Error = json.Unmarshal(request.Body, &args); response.Error == nil {
				result := make(chan *daemon.Process)
				if response.Error = dm.StartProgram(args[0], result); response.Error == nil {
					select {
					case p := <-result:
						if p.GetStatus() == daemon.Running {
							response.Body = OK
						} else {
							response.Error = errors.New("start error:" + p.GetStatus().String())
						}
					case <-time.After(time.Second * 7):
						response.Body = []byte("wait status timeout")
					}
				}
				commons.SafeExec(func() { close(result) })
			}
		case "stop":
			args := []string{}
			if response.Error = json.Unmarshal(request.Body, &args); response.Error == nil {
				if response.Error = dm.StopProgram(args[0]); response.Error == nil {
					response.Body = OK
				}
			}
		case "detail":
			args := []string{}
			if response.Error = json.Unmarshal(request.Body, &args); response.Error == nil {
				if pro, err := dm.GetProgram(args[0]); err != nil {
					response.Error = err
				} else {
					pro.Refresh()
					response.Body, _ = json.Marshal(pro)
				}
			}
		case "status":
			var pro *daemon.Process
			args := []string{}
			if response.Error = json.Unmarshal(request.Body, &args); response.Error == nil {
				if pro, response.Error = dm.GetProgram(args[0]); response.Error == nil {
					response.Body = []byte(pro.GetStatus())
				}
			}
		case "list":
			args := []string{}
			if response.Error = json.Unmarshal(request.Body, &args); response.Error == nil {
				if len(args) == 0 {
					names := dm.ListProgramNames()
					response.Body, _ = json.Marshal(names)
				} else if args[0] == "inspect" {
					process := dm.ListProcess()
					for _, p := range process {
						p.Refresh()
					}
					response.Body, _ = json.Marshal(process)
				} else {
					response.Error = errors.New("Invalid parameter：" + strings.Join(args, " "))
				}
			}
		case "add":
			{
				args := []string{}
				if response.Error = json.Unmarshal(request.Body, &args); response.Error == nil {
					body := "result:"
					for _, arg := range args {

						program := daemon.NewProgram()
						if err := json.Unmarshal([]byte(arg), program); err != nil {
							response.Error = errors.New(fmt.Sprintf("\nadd program %s error: %s", program.Name, err))
							return response
						}
						program.Id = dm.MaxId()
						if err := dm.AddProgram(program); err != nil {
							response.Error = errors.New(fmt.Sprintf("\nadd program %s error: %s", program.Name, err))
							return response
						}
						if response.Error = dm.WriteConfig(program.Name); response.Error == nil {
							body += fmt.Sprintf("\nadd program %s ok", program.Name)
						}
					}
					response.Body = []byte(body)
				}
			}
		case "delete":
			{
				args := []string{}
				if response.Error = json.Unmarshal(request.Body, &args); response.Error == nil {
					skipHeader, _ := request.GetHeader("skip")
					skip, _ := strconv.ParseBool(skipHeader)
					body := "result:"
					for _, arg := range args {
						if err := dm.RemoveProgram(arg, skip); err != nil {
							response.Error = errors.New(fmt.Sprintf("\nremove program %s error: %s", arg, err))
							return response
						} else {
							body += fmt.Sprintf("\nremove program %s ok", arg)
							if err = files.New(conf.Config.Server.Dir + "/" + arg + ".json").Remove(); err != nil {
								body += fmt.Sprintf("\nremove program %s config file error: %s", arg, err)
							}
						}
					}
					response.Body = []byte(body)
				}
			}
		case "modify":
			{
				program := daemon.NewProgram()
				if response.Error = json.Unmarshal(request.Body, program); response.Error == nil {
					if p, err := dm.GetProgram(program.Name); err != nil {
						response.Error = err
					} else if p.GetStatus().IsRunning() {
						response.Error = errors.New("You can not edit a running program")
					} else {
						if response.Error = dm.ModifyProgram(program); response.Error == nil {
							if response.Error = dm.WriteConfig(program.Name); response.Error == nil {
								response.Body = OK
							}
						}
					}
				}
			}
		default:
			response.Error = errors.New("InvalidCommand")
		}
		return response
	}
}
