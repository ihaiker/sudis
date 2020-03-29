package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/daemon"
	"github.com/ihaiker/sudis/libs/config"
	. "github.com/ihaiker/sudis/libs/errors"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/join"
	"strconv"
	"strings"
	"time"
)

var logger = logs.GetLogger("server")

func getNode(request *rpc.Request) string {
	node, has := request.GetHeader("node")
	if !has {
		node = config.Config.Key
	}
	return node
}

func getArgs(request *rpc.Request) (args []string, name string) {
	args = make([]string, 0)
	errors.Assert(json.Unmarshal(request.Body, &args))
	if len(args) > 0 {
		name = args[0]
	}
	return
}

func getTimeout(request *rpc.Request) time.Duration {
	timeout, has := request.GetHeader("timeout")
	if !has {
		return time.Second * 7
	}
	if seconds, err := strconv.Atoi(timeout); err != nil {
		return time.Second * 7
	} else {
		return time.Second * time.Duration(seconds)
	}
}

func True(request *rpc.Request, header string) bool {
	val, exists := request.GetHeader(header)
	return exists && val == "true"
}

func MakeCommand(dm *cluster.DaemonManager, joinManager *join.ToJoinManager) rpc.OnMessage {

	return func(channel remoting.Channel, request *rpc.Request) (response *rpc.Response) {
		response = rpc.NewResponse(request.ID())
		node := getNode(request)
		timeout := getTimeout(request)

		logger.Debugf("node: %s, action: %s, headers: %v, args: %s",
			node, request.URL, request.Headers, string(request.Body))

		defer errors.Catch(func(err error) {
			if _, match := err.(*errors.WrapError); !match {
				logger.Debug("cli handler error: ", err)
			}
			response.Error = errors.Root(err)
		})

		switch request.URL {
		case "start", "stop", "delete":
			names, _ := getArgs(request)
			body := bytes.NewBufferString("")
			for _, name := range names {
				err := dm.Command(node, name, request.URL, timeout)
				if err != nil {
					body.WriteString(fmt.Sprintf("%s %s : %s\n", name, request.URL, err))
				} else {
					body.WriteString(fmt.Sprintf("%s %s : OK\n", name, request.URL))
				}
			}
			response.Body = body.Bytes()

		case "detail":
			_, name := getArgs(request)
			pro := dm.MustGetProcess(node, name)
			response.Body, _ = json.Marshal(pro)

		case "status":
			_, name := getArgs(request)
			pro := dm.MustGetProcess(node, name)
			response.Body = []byte(pro.GetStatus())

		case "list":
			node, _ = request.GetHeader("node")

			all := True(request, "all")
			inspect := True(request, "inspect")
			quiet := True(request, "quiet")

			processes := dm.ListPrograms("", node, "", "", 1, 2000)

			outs := make([]interface{}, 0)
			if !inspect && !quiet {
				outs = append(outs, strings.Repeat("-", 2+20+3+20+3+7+2))
				outs = append(outs, fmt.Sprintf("| %20s | %20s | %7s |", "Node", "Name", "Status"))
				outs = append(outs, strings.Repeat("-", 2+20+3+20+3+7+2))
			}

			for _, p := range processes.Data.([]*daemon.Process) {
				if all || p.Status.IsRunning() {
					if inspect {
						outs = append(outs, p)
					} else if quiet {
						outs = append(outs, fmt.Sprintf("%s.%s", p.Node, p.Name))
					} else {
						outs = append(outs, fmt.Sprintf("| %20s | %20s | %7s |", p.Node, p.Name, p.Status.String()))
						outs = append(outs, strings.Repeat("-", 2+20+3+20+3+7+2))
					}
				}
			}

			response.Body, _ = json.Marshal(outs)

		case "add":
			{
				args, _ := getArgs(request)
				body := bytes.NewBufferString("")
				for _, arg := range args {
					program := daemon.NewProgram()
					if err := json.Unmarshal([]byte(arg), program); err != nil {
						body.WriteString(fmt.Sprintf("\nadd program %s error: %s", program.Name, err))
					} else if err := dm.AddProgram(node, program); err != nil {
						body.WriteString(fmt.Sprintf("\nadd program %s error: %s", program.Name, err))
					} else {
						body.WriteString(fmt.Sprintf("\nadd program %s ok", program.Name))
					}
				}
				response.Body = body.Bytes()
			}
		case "modify":
			{
				program := daemon.NewProgram()
				errors.Assert(json.Unmarshal(request.Body, program))
				name := program.Name
				p := dm.MustGetProcess(node, program.Name)
				if p.GetStatus().IsRunning() {
					response.Error = ErrProgramIsRunning
					return
				} else {
					dm.MustModifyProgram(node, name, program)
				}
				response.Body = []byte("OK")
			}
		case "tail":
			{
				args, name := getArgs(request)
				name, reg, id := args[0], args[1], args[2]
				line := 30
				if numStr, has := request.GetHeader("num"); has {
					line, _ = strconv.Atoi(numStr)
				}
				logger.Debugf("接收日志消息。name: %s, reg: %s, id: %s, num: %d", name, reg, id, line)
				if reg == "true" {
					response.Error = dm.SubscribeLogger(id, node, name, func(id, line string) {
						_ = errors.Safe(func() error {
							req := new(rpc.Request)
							req.URL = "tail.logger"
							req.Header("id", id)
							req.Body = []byte(line)
							return channel.Write(req, time.Second*3)
						})
					}, line)
				} else {
					response.Error = dm.UnsubscribeLogger(id, node, name)
				}
			}
		case "tag":
			{
				args, name := getArgs(request)
				tag := args[1]
				delete, has := request.GetHeader("delete")
				add := !has || delete == "false"
				errors.Assert(dm.ModifyProgramTag(name, node, tag, add), "modify tag")
				response.Body = []byte("OK")
			}
		case "join":
			{
				addresses, _ := getArgs(request)
				must, has := request.GetHeader("must")
				token, exists := request.GetHeader("token")
				if token == "" || !exists {
					token = config.Config.Salt
				}

				if token == "" {
					response.Error = ErrToken
					return
				}

				out := bytes.NewBufferString("")
				var err error
				for _, address := range addresses {
					if has && must == "true" {
						joinManager.MustJoinIt(address, token)
						err = nil
					} else {
						err = joinManager.Join(address, token)
					}
					if err != nil {
						out.WriteString(fmt.Sprintf("join %s %s", address, err.Error()))
					} else {
						out.WriteString(fmt.Sprintf("join %s OK", address))
					}
				}
				response.Body = out.Bytes()
			}
		case "leave":
			{
				address, _ := getArgs(request)
				if err := joinManager.Leave(address...); err != nil {
					response.Body = []byte(fmt.Sprintf("leave: %s", err))
				} else {
					response.Body = []byte("OK")
				}
			}
		default:
			response.Error = errors.New("InvalidCommand")
		}
		return
	}
}
