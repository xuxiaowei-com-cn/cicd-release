package api

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/cicd-release/constant"
	"log"
	"os"
)

func Gitee(prerelease bool, context *cli.Context) error {
	var tag = context.String(constant.Tag)
	var autoCreateTag = context.Bool(constant.AutoCreateTag)
	var gitlabExportAssetsFileName = context.String(constant.GitlabExportAssetsFileName)

	log.Printf("是否是预发布版本：%v", prerelease)
	log.Printf("发布到 Gitee")

	err := AutoCreateTag(tag, autoCreateTag)
	if err != nil {
		return err
	}

	if gitlabExportAssetsFileName != "" {
		jsonData, err := os.ReadFile(gitlabExportAssetsFileName)
		if err != nil {
			log.Printf("ReadFile %s Error:\n%s", gitlabExportAssetsFileName, err)
			return err
		}

		readResult := make(map[string]interface{})

		err = json.Unmarshal(jsonData, &readResult)
		if err != nil {
			log.Printf("Unmarshal %s Error:\n%s", gitlabExportAssetsFileName, err)
			return err
		}

		fmt.Println(readResult)
	}

	return nil
}
