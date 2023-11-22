package command

import (
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/cicd-release/api"
	"github.com/xuxiaowei-com-cn/cicd-release/flag"
)

func ReleaseCommand() *cli.Command {
	return &cli.Command{
		Name:  "release",
		Usage: "发布",
		Flags: flag.CommonFlag(),
		Before: func(context *cli.Context) error {
			_, err := api.GitVersion()

			return err
		},
		Subcommands: []*cli.Command{
			{
				Name:  "gitee",
				Usage: "Gitee 发布",
				Flags: flag.GiteeFlag(),
				Action: func(context *cli.Context) error {

					return api.Gitee(false, context)
				},
			},
			{
				Name:  "gitlab",
				Usage: "GitLab 发布",
				Flags: flag.GitlabFlag(),
				Action: func(context *cli.Context) error {

					return api.Gitlab(false, context)
				},
			},
			{
				Name:  "github",
				Usage: "GitHub 发布",
				Flags: flag.GithubFlag(),
				Action: func(context *cli.Context) error {

					return api.Github(false, context)
				},
			},
		},
	}
}
