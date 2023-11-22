package api

import (
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/cicd-release/constant"
	"log"
)

func Gitlab(prerelease bool, context *cli.Context) error {
	var gitlabInstance = context.String(constant.GitlabInstance)

	log.Printf("是否是预发布版本：%v", prerelease)
	log.Printf("发布到 GitLab，实例：%s", gitlabInstance)

	err := AutoCreateTag(tag, autoCreateTag)
	if err != nil {
		return err
	}

	return nil
}
