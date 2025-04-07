package utils

import (
	"bytes"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
)

func GetRootGolangModule() (string, error) {
	// 获取根模块名
	cmdOutput, err := exec.Command("go", "list", "-m").CombinedOutput()
	if err != nil {
		return "", err
	}

	// 按换行符拆分结果
	lines := bytes.Split(cmdOutput, []byte("\n"))
	if len(lines) == 0 {
		return "", errors.New("source code may not using go module")
	}

	return strings.TrimSpace(string(lines[0])), nil
}

func Fatal(prompt string, errs ...error) {
	if len(errs) > 0 {
		fmt.Printf("%s: %s\n", prompt, errs[0].Error())
	} else {
		fmt.Println(prompt)
	}
	os.Exit(1)
}

// GetInput 获取字符串输入
func GetInput(prompt string) string {
	rl, _ := readline.New(prompt)
	defer func() {
		if rl != nil {
			rl.Close()
		}
	}()

	for {
		line, _ := rl.Readline()
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
		fmt.Println("输入不能为空，请重新输入")
	}
}
