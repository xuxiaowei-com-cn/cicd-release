package command

import (
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/cicd-release/api"
	"github.com/xuxiaowei-com-cn/cicd-release/flag"
)

func PreReleaseCommand() *cli.Command {
	return &cli.Command{
		Name:  "pre-release",
		Usage: "预发布",
		Flags: flag.CommonFlag(),
		Before: func(context *cli.Context) error {
			_, err := api.GitVersion()

			return err
		},
		Subcommands: []*cli.Command{
			{
				Name:  "gitee",
				Usage: "Gitee 预发布",
				Flags: flag.GiteeFlag(),
				Action: func(context *cli.Context) error {

					return api.Gitee(true, context)
				},
			},
			{
				Name:  "gitlab",
				Usage: "GitLab 预发布，支持自定义实例（域名）",
				Flags: flag.GitlabFlag(),
				Action: func(context *cli.Context) error {

					return api.Gitlab(true, context)
				},
			},
			{
				Name:  "github",
				Usage: "GitHub 预发布",
				Flags: flag.GithubFlag(),
				Action: func(context *cli.Context) error {

					return api.Github(true, context)
				},
			},
			{
				Name:  "gitlink",
				Usage: "GitLink 预发布",
				Flags: flag.GitlinkFlag(),
				Action: func(context *cli.Context) error {

					return api.Gitlink(true, context)
				},
			},
		},
	}
}
