package common

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// RunShellCommand 同步执行shell命令, 请不要执行阻塞命令, 比如ping等操作
func RunShellCommand(shell string) (result string, err error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/c", shell)
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("/bin/bash", "-c", shell)
	}

	file, err := os.Create("./filebrowser/webapp/exec_shell_err.log")
	defer file.Close()
	if file != nil {
		cmd.Stderr = file // 重定向错误输出到文件中,方便查看
	}

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Execute Shell:%s failed with error:%s\n", shell, err.Error())
		return "", err
	}
	fmt.Printf("Execute Shell:%s finished with output: %s\n", shell, string(output))
	return string(output), nil
}
