package actions

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"syscall"

	"github.com/caddyserver/caddy/v2/web/cache"
	"github.com/caddyserver/caddy/v2/web/common"
	"github.com/caddyserver/caddy/v2/web/model"
	"github.com/kataras/iris/v12"
)

type execCancelContext struct {
	cancel  func()
	process *os.Process
}

var cancelDomainFuncMaps map[string]execCancelContext
var lock sync.Mutex

func init() {
	cancelDomainFuncMaps = make(map[string]execCancelContext)
}

// ExecShell 是执行shell的api方法
func ExecShell(ctx iris.Context) {
	var domain = ctx.FormValue("domain")
	var shellType = ctx.FormValue("shell_type")

	var shellConfigMap map[string]model.ShellConfig
	shellConfigMap = make(map[string]model.ShellConfig)
	err := cache.GetCacheData("SyncShellConfigList", &shellConfigMap)
	if err != nil {
		ctx.JSON(model.ResponseData{State: false, Message: "读取配置失败", Error: err.Error(), HTTPCode: 500})
		return
	}

	if shellConfig, ok := shellConfigMap[domain]; ok {
		var shell = ""
		if runtime.GOOS == "windows" {
			shell = fmt.Sprintf("cd ./filebrowser/webapp/%s & ", shellConfig.Domain)
		} else if runtime.GOOS == "linux" {
			shell = fmt.Sprintf("cd ./filebrowser/webapp/%s && ", shellConfig.Domain)
		}

		var result = ""
		switch shellType {
		case "start":
			if len(shellConfig.StartShell) > 0 {
				if _, ok := cancelDomainFuncMaps[domain]; ok {
					result = "process is already startd."
				} else {
					resultChan := make(chan string)
					go func() {
						var cmd *exec.Cmd
						var sysProcAttr *syscall.SysProcAttr

						ctx, cancel := context.WithCancel(context.Background())
						if runtime.GOOS == "windows" {
							cmd = exec.CommandContext(ctx, "cmd.exe", "/c", shell+shellConfig.StartShell)
						} else if runtime.GOOS == "linux" {
							fmt.Println("执行的shell: " + shell + shellConfig.StartShell)
							cmd = exec.CommandContext(ctx, "/bin/bash", "-c", shell+shellConfig.StartShell)
						}

						file, err := os.Create("./filebrowser/webapp/" + domain + "/exec_shell_err.log")
						defer file.Close()

						sysProcAttr = &syscall.SysProcAttr{
							Setpgid: true, // 使子进程拥有自己的 pgid，等同于子进程的 pid
						}

						cmd.SysProcAttr = sysProcAttr
						cmd.Stdout = os.Stdout
						// cmd.Stderr = file

						err = cmd.Start()
						if err != nil {
							fmt.Printf("\n cmd.Start err: %v",err)
						}

						lock.Lock()
						execContext := execCancelContext{}
						execContext.cancel = cancel
						execContext.process = cmd.Process
						cancelDomainFuncMaps[domain] = execContext
						lock.Unlock()

						if cmd.Process != nil {
							resultChan <- "start process pid:" + strconv.Itoa(cmd.Process.Pid)
						} else {
							resultChan <- "start process pid: nil"
						}
						close(resultChan)
						/*
							err = cmd.Process.Kill()
							if err != nil {
								fmt.Printf("process kill err: %v\n", err)
							}
						*/
						err = cmd.Wait()
						if err != nil {
							fmt.Printf("\nwait for exec err: %v", err)
						}

						lock.Lock()
						delete(cancelDomainFuncMaps, domain)
						lock.Unlock()
					}()
					result = <-resultChan
				}
			}
		case "stop":
			if cancelContext, ok := cancelDomainFuncMaps[domain]; ok {
				syscall.Kill(-cancelContext.process.Pid, syscall.SIGKILL)
				cancelContext.cancel()
				err = cancelContext.process.Kill()
				if err != nil {
					fmt.Printf("stop exec err: %v\n", err)
				}

				lock.Lock()
				delete(cancelDomainFuncMaps, domain)
				lock.Unlock()
			}
		case "sync":
			if len(shellConfig.SyncShell) > 0 {
				result, err = common.RunShellCommand(shell + shellConfig.SyncShell)
			}
		default:
			ctx.JSON(model.ResponseData{State: false, Message: "无效的指令:" + shellType, HTTPCode: 400})
			return
		}
		fmt.Printf("\nExecShell执行结果, result: %s , err: %s \n", result, err)
		ctx.JSON(model.ResponseData{State: true, Message: "shell执行完成", Data: result, HTTPCode: 400})
		return
	}

	ctx.JSON(model.ResponseData{State: false, Message: "无效的参数domain:" + domain, HTTPCode: 400})
	return

}
