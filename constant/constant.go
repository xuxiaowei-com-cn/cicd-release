package constant

const (
	PreRelease                 = "pre-release"                    // 预发布
	Release                    = "release"                        // 发布
	ReleaseName                = "release-name"                   // 发布名称
	ReleaseBody                = "release-body"                   // 发布详情
	Tag                        = "tag"                            // 发布标签
	PackageName                = "package-name"                   // 包名，只能包含小写字母（az）、大写字母（AZ）、数字（0-9）、点（.）、连字符（-）或下划线（_）
	Milestones                 = "milestones"                     // 发布里程碑
	AutoCreateTag              = "auto-create-tag"                // 自动创建不存在的标签
	Artifacts                  = "artifacts"                      // 发布产物
	GithubToken                = "github-token"                   // Github Token
	GithubRepository           = "github-repository"              // Github 仓库，如：https://github.com/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release
	GiteeUsername              = "gitee-username"                 // Gitee 用户名
	GiteeToken                 = "gitee-token"                    // Gitee Token
	GiteeRepository            = "gitee-repository"               // Gitee 仓库，如：https://gitee.com/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release
	GitlabToken                = "gitlab-token"                   // Gitlab Token
	GitlabInstance             = "gitlab-instance"                // Gitlab 实例（协议 + 域名）
	GitlabApi                  = "gitlab-api"                     // Gitlab API
	GitlabRepository           = "gitlab-repository"              // Gitee 仓库，如：https://gitlab.com/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release
	GitlabExportAssetsFileName = "gitlab-export-assets-file-name" // Gitlab 导出资源文件名称
)
