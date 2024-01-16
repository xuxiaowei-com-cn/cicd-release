package constant

const (
	PreRelease                  = "pre-release"                     // 预发布
	Release                     = "release"                         // 发布
	ReleaseName                 = "release-name"                    // 发布名称
	ReleaseBody                 = "release-body"                    // 发布详情
	Tag                         = "tag"                             // 发布标签
	Draft                       = "draft"                           // 草稿
	PackageName                 = "package-name"                    // 包名，只能包含小写字母（az）、大写字母（AZ）、数字（0-9）、点（.）、连字符（-）或下划线（_）
	Milestones                  = "milestones"                      // 发布里程碑
	AutoCreateTag               = "auto-create-tag"                 // 自动创建不存在的标签
	Artifacts                   = "artifacts"                       // 发布产物
	GithubUsername              = "github-username"                 // Github 用户名
	GithubToken                 = "github-token"                    // Github Token
	GithubRepository            = "github-repository"               // Github 仓库，如：https://github.com/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release
	GiteeUsername               = "gitee-username"                  // Gitee 用户名
	GiteeToken                  = "gitee-token"                     // Gitee Token
	GiteeRepository             = "gitee-repository"                // Gitee 仓库，如：https://gitee.com/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release
	GitlabUsername              = "gitlab-username"                 // Gitlab 用户名
	GitlabToken                 = "gitlab-token"                    // Gitlab Token
	GitlabInstance              = "gitlab-instance"                 // Gitlab 实例（协议 + 域名）
	GitlabApi                   = "gitlab-api"                      // Gitlab API
	GitlabRepository            = "gitlab-repository"               // Gitee 仓库，如：https://gitlab.com/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release
	GitlabExportAssetsFileName  = "gitlab-export-assets-file-name"  // Gitlab 导出资源文件名称
	GitlabImportAssetsFileName  = "gitlab-import-assets-file-name"  // Gitlab 导入资源文件名称
	GitlinkUsername             = "gitlink-username"                // gitlink 用户名
	GitlinkToken                = "gitlink-token"                   // gitlink Token
	GitlinkCookie               = "gitlink-cookie"                  // gitlink Cookie
	GitlinkRepository           = "gitlink-repository"              // gitlink 仓库，如：https://gitlink.org.cn/xuxiaowei-com-cn/cicd-release.git 仓库应该为：xuxiaowei-com-cn/cicd-release
	GitlinkExportAssetsFileName = "gitlink-export-assets-file-name" // gitlink 导出资源文件名称
	GitlinkAttachmentsPrefix    = "gitlink-attachments-prefix"      // gitlink 附件URL前缀，如：https://www.gitlink.org.cn/api/attachments
)
