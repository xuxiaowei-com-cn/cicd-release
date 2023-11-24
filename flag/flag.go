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

func ReleaseBodyFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.ReleaseBody,
		Usage:    "发布详情",
		Required: required,
	}
}

func TagFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.Tag,
		Usage:    "发布标签",
		Required: required,
	}
}

func DraftFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:  constant.Draft,
		Usage: "Github 草稿",
		Value: false,
	}
}

func PackageNameFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.PackageName,
		Usage:    "包名，即：GitLab 产物储存 URL 前缀。\n\t只能包含小写字母（az）、大写字母（AZ）、数字（0-9）、点（.）、连字符（-）或下划线（_）",
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
		Usage: "发布产物（包含路径）。\n\t可以包含多级路径。\n\t文件名（除路径外，所有文件名均不能出现重复）：只能包含小写字母（az）、大写字母（AZ）、数字（0-9）、点（.）、连字符（-）或下划线（_）。",
	}
}

func GithubUsernameFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.GithubUsername,
		Usage:    "Github 用户名",
		EnvVars:  []string{"GITHUB_ACTOR"},
		Required: required,
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
		Usage:    "Github 仓库。\n\t如：https://github.com/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release",
		EnvVars:  []string{"GITHUB_REPOSITORY"},
		Required: required,
	}
}

func GiteeUsername(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.GiteeUsername,
		Usage:    "Gitee 用户名",
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
		Usage:    "Gitee 仓库。\n\t如：https://gitee.com/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release",
		EnvVars:  []string{"GITEE_REPO"},
		Required: required,
	}
}

func GitlabUsernameFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     constant.GitlabUsername,
		Usage:    "Gitlab 用户名",
		EnvVars:  []string{"GITLAB_USER_LOGIN"},
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
		Usage:    "Gitlab 仓库。\n\t如：https://gitlab.com/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release",
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
		Usage: "Gitlab API",
		Value: "api/v4",
	}
}

func GitlabExportAssetsNameFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  constant.GitlabExportAssetsFileName,
		Usage: "Gitlab 导出资源文件名称。\n\t主要用于发布到 Gitee 时在版本发布中新增产物下载地址（Gitee 没有上传产物的 API）。\n\t导出格式为 map，键：代表文件名，值：代表下载链接",
	}
}

func GiteeFlag() []cli.Flag {
	return []cli.Flag{
		ReleaseNameFlag(true),
		ReleaseBodyFlag(true),
		TagFlag(true),
		AutoCreateTagFlag(),
		GitlabExportAssetsNameFlag(),

		GiteeRepositoryFlag(true),
		GiteeUsername(true),
		GiteeTokenFlag(true),
	}
}

func GitlabFlag() []cli.Flag {
	return []cli.Flag{
		ReleaseNameFlag(true),
		ReleaseBodyFlag(false),
		TagFlag(true),
		PackageNameFlag(true),
		AutoCreateTagFlag(),
		MilestonesFlag(),
		ArtifactsFlag(),

		GitlabInstanceFlag(),
		GitlabApiFlag(),
		GitlabRepositoryFlag(true),
		GitlabUsernameFlag(true),
		GitlabTokenFlag(true),
		GitlabExportAssetsNameFlag(),
	}
}

func GithubFlag() []cli.Flag {
	return []cli.Flag{
		ReleaseNameFlag(true),
		ReleaseBodyFlag(false),
		TagFlag(true),
		DraftFlag(),
		AutoCreateTagFlag(),
		ArtifactsFlag(),

		GithubRepositoryFlag(true),
		GithubUsernameFlag(true),
		GithubTokenFlag(true),
	}
}

func CommonFlag() []cli.Flag {
	return []cli.Flag{
		ReleaseNameFlag(false),
		ReleaseBodyFlag(false),
		TagFlag(false),
		DraftFlag(),
		PackageNameFlag(false),
		AutoCreateTagFlag(),
		MilestonesFlag(),
		ArtifactsFlag(),

		GiteeRepositoryFlag(false),
		GiteeUsername(false),
		GiteeTokenFlag(false),

		GitlabInstanceFlag(),
		GitlabApiFlag(),
		GitlabRepositoryFlag(false),
		GitlabUsernameFlag(false),
		GitlabTokenFlag(false),
		GitlabExportAssetsNameFlag(),

		GithubRepositoryFlag(false),
		GithubUsernameFlag(false),
		GithubTokenFlag(false),
	}
}
