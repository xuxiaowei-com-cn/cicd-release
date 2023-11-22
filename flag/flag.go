package flag

import (
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/cicd-release/constant"
)

func ReleaseNameFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.ReleaseName,
		Usage:    "发布名称",
		Required: required,
	}
}

func ReleaseBodyFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  constant.ReleaseBody,
		Usage: "发布详情",
	}
}

func TagFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.Tag,
		Usage:    "发布标签",
		Required: required,
	}
}

func AutoCreateTagFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:  constant.AutoCreateTag,
		Usage: "是否自动创建不存在的标签",
		Value: false,
	}
}

func MilestonesFlag() cli.Flag {
	return &cli.StringSliceFlag{
		Name:  constant.Milestones,
		Usage: "发布里程碑",
	}
}

func ArtifactsFlag() cli.Flag {
	return &cli.StringSliceFlag{
		Name:  constant.Artifacts,
		Usage: "发布产物",
	}
}

func GithubTokenFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.GithubToken,
		Usage:    "Github Token",
		EnvVars:  []string{"GITHUB_TOKEN"},
		Required: required,
	}
}

func GithubRepositoryFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.GithubRepository,
		Usage:    "Github 仓库，如：https://github.com/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release",
		EnvVars:  []string{"GITHUB_REPOSITORY"},
		Required: required,
	}
}

func GiteeTokenFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.GiteeToken,
		Usage:    "Gitee Token",
		Required: required,
	}
}

func GiteeRepositoryFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.GiteeRepository,
		Usage:    "Gitee 仓库，如：https://gitee.com/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release",
		EnvVars:  []string{"GITEE_REPO"},
		Required: required,
	}
}

func GitlabTokenFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.GitlabToken,
		Usage:    "Gitlab Token",
		Required: required,
	}
}

func GitlabRepositoryFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.GitlabRepository,
		Usage:    "Gitlab 仓库，如：https://gitlab.com/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release",
		EnvVars:  []string{"CI_PROJECT_PATH"},
		Required: required,
	}
}

func GitlabInstanceFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    constant.GitlabInstance,
		Usage:   "Gitlab 实例（协议 + 域名）",
		Value:   "https://gitlab.com",
		EnvVars: []string{"CI_SERVER_URL"},
	}
}

func GitlabApiFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  constant.GitlabApi,
		Usage: "Gitlab API，如：api/v4",
		Value: "api/v4",
	}
}

func GiteeFlag() []cli.Flag {
	return []cli.Flag{
		ReleaseNameFlag(true),
		ReleaseBodyFlag(),
		TagFlag(true),
		AutoCreateTagFlag(),
		MilestonesFlag(),
		ArtifactsFlag(),

		GiteeRepositoryFlag(true),
		GiteeTokenFlag(true),
	}
}

func GitlabFlag() []cli.Flag {
	return []cli.Flag{
		ReleaseNameFlag(true),
		ReleaseBodyFlag(),
		TagFlag(true),
		AutoCreateTagFlag(),
		MilestonesFlag(),
		ArtifactsFlag(),

		GitlabInstanceFlag(),
		GitlabApiFlag(),
		GitlabRepositoryFlag(true),
		GitlabTokenFlag(true),
	}
}

func GithubFlag() []cli.Flag {
	return []cli.Flag{
		ReleaseNameFlag(true),
		ReleaseBodyFlag(),
		TagFlag(true),
		AutoCreateTagFlag(),
		MilestonesFlag(),
		ArtifactsFlag(),

		GithubRepositoryFlag(true),
		GithubTokenFlag(true),
	}
}

func CommonFlag() []cli.Flag {
	return []cli.Flag{
		ReleaseNameFlag(false),
		ReleaseBodyFlag(),
		TagFlag(false),
		AutoCreateTagFlag(),
		MilestonesFlag(),
		ArtifactsFlag(),

		GiteeRepositoryFlag(false),
		GiteeTokenFlag(false),

		GitlabInstanceFlag(),
		GitlabApiFlag(),
		GitlabRepositoryFlag(false),
		GitlabTokenFlag(false),

		GithubRepositoryFlag(false),
		GithubTokenFlag(false),
	}
}
