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
func GetInput(prompt string, defaults ...string) string {
	rlConfig := &readline.Config{}

	var defaultValue string
	if len(defaults) > 0 {
		defaultValue = strings.TrimSpace(defaults[0])
	}

	if defaultValue != "" {
		rlConfig.Prompt = fmt.Sprintf("%s[%s]: ", prompt, defaultValue)
	} else {
		rlConfig.Prompt = fmt.Sprintf("%s: ", prompt)
	}

	rl, _ := readline.NewEx(rlConfig)
	defer func() {
		if rl != nil {
			rl.Close()
		}
	}()

	var inputValue string
	for {
		line, err := rl.Readline()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				os.Exit(0)
			}
			
			if defaultValue != "" {
				inputValue = defaultValue
			}
			break
		}

		line = strings.TrimSpace(line)
		if line != "" {
			inputValue = line
		} else if defaultValue != "" {
			inputValue = defaultValue
		}

		if inputValue != "" {
			break
		}

		fmt.Println("empty input, please input again!")
	}

	return inputValue
}
