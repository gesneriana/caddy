package common

import (
	"fmt"
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

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Execute Shell:%s failed with error:%s\n", shell, err.Error())
		return "", err
	}
	fmt.Printf("Execute Shell:%s finished with output: %s\n", shell, string(output))
	return string(output), nil
}
