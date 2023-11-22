package api

import (
	"fmt"
	"log"
	"os/exec"
)

// GitVersion 获取 Git 版本号
func GitVersion() (string, error) {
	cmd := exec.Command("git1", "--version")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("无法获取 Git 版本号")
		return "", err
	}

	version := string(output)
	fmt.Println("Git 版本:", version)
	return version, nil
}
