package api

import (
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/cicd-release/constant"
	"log"
	"os/exec"
)

// GitVersion 获取 Git 版本号
func GitVersion() (string, error) {
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("无法获取 Git 版本号")
		return "", err
	}

	version := string(output)
	log.Println("Git 版本:", version)
	return version, nil
}

func GitPrintTag(tag string) error {
	cmd := exec.Command("git", "show", tag, "--no-patch")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Git 标签 %s 不存在\n", tag)
		return err
	}

	log.Printf("Git 标签 %s 信息: \n%s\n", tag, output)
	return nil
}

func GitCreateTag(tag string) error {
	cmd := exec.Command("git", "tag", tag)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("Git 创建标签 %s 异常\n", tag)
		return err
	}

	return GitPrintTag(tag)
}

func AutoCreateTag(context *cli.Context) error {
	var tag = context.String(constant.Tag)
	var autoCreateTag = context.Bool(constant.AutoCreateTag)

	err := GitPrintTag(tag)
	if err != nil {
		if autoCreateTag {
			log.Printf("开始 创建标签：%s\n", tag)
			err = GitCreateTag(tag)
			if err != nil {
				return err
			}
			log.Printf("完成 创建标签：%s\n", tag)

			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
