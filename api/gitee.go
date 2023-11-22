package api

import (
	"github.com/urfave/cli/v2"
	"log"
)

func Gitee(prerelease bool, context *cli.Context) error {
	log.Printf("是否是预发布版本：%v", prerelease)
	log.Printf("发布到 Gitee")

	err := AutoCreateTag(context)
	if err != nil {
		return err
	}

	return nil
}
