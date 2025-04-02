package utils

import (
	"bufio"
	"bytes"
	"fmt"
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

func Fatal(prompt string, err error) {
	fmt.Printf("%s: %s\n", prompt, err.Error())
	os.Exit(1)
}

// GetInput 获取字符串输入
func GetInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			return input
		}
		fmt.Println("输入不能为空，请重新输入")
	}
}
