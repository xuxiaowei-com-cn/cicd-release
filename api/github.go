package api

import (
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/cicd-release/constant"
	"log"
)

func Github(prerelease bool, context *cli.Context) error {
	var tag = context.String(constant.Tag)
	var autoCreateTag = context.Bool(constant.AutoCreateTag)

	log.Printf("是否是预发布版本：%v", prerelease)
	log.Printf("发布到 GitHub")

	err := AutoCreateTag(tag, autoCreateTag)
	if err != nil {
		return err
	}

	return nil
}
