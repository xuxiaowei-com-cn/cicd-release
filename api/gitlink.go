package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/cicd-release/constant"
	"github.com/xuxiaowei-com-cn/go-gitlink/v2"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Gitlink(prerelease bool, context *cli.Context) error {

	var releaseName = context.String(constant.ReleaseName)
	var releaseBody = context.String(constant.ReleaseBody)
	var tag = context.String(constant.Tag)
	var autoCreateTag = context.Bool(constant.AutoCreateTag)
	var artifacts = context.StringSlice(constant.Artifacts)
	var gitlinkRepository = context.String(constant.GitlinkRepository)
	var gitlinkUsername = context.String(constant.GitlinkUsername)
	var gitlinkToken = context.String(constant.GitlinkToken)
	var gitlinkCookie = context.String(constant.GitlinkCookie)
	var draft = context.Bool(constant.Draft)
	var gitlinkExportAssetsFileName = context.String(constant.GitlinkExportAssetsFileName)
	var gitlinkAttachmentsPrefix = context.String(constant.GitlinkAttachmentsPrefix)

	log.Printf("是否是预发布版本：%v", prerelease)
	log.Printf("发布到 Gitlink，路径：%s", gitlinkRepository)

	_, err := url.Parse(gitlinkAttachmentsPrefix)
	if err != nil {
		return err
	}

	// 检查发布
	err = GitlinkGetReleases(gitlinkCookie)
	if err != nil {
		return err
	}

	// 自动创建标签
	err = AutoCreateTag(tag, autoCreateTag)
	if err != nil {
		return err
	}

	// 检查标签
	err = GitlinkGetTag(gitlinkCookie, tag)
	if err != nil {
		return err
	}

	// 推送标签
	err = GitPushTag("https://gitlink.org.cn", gitlinkRepository, gitlinkUsername, gitlinkToken, tag)
	if err != nil {
		return err
	}

	// 上传产物
	attachmentIds, _, err := GitlinkAttachments(artifacts, gitlinkExportAssetsFileName, gitlinkCookie)
	if err != nil {
		return err
	}

	// 发布
	err = GitlinkReleases(prerelease, releaseName, releaseBody, tag, gitlinkRepository, gitlinkCookie, attachmentIds, draft)
	if err != nil {
		return err
	}

	return nil
}

// GitlinkGetReleases
// 检查发布
func GitlinkGetReleases(gitlinkCookie string) error {

	return nil
}

// GitlinkGetTag
// 检查标签
func GitlinkGetTag(gitlinkCookie string, tag string) error {

	return nil
}

// GitlinkAttachments
// 上传产物
func GitlinkAttachments(artifacts []string, gitlinkExportAssetsFileName string, gitlinkCookie string) ([]string, map[string]interface{}, error) {

	gitClient, err := gitlink.NewClient("")
	if err != nil {
		return nil, nil, err
	}

	gitClient.Cookie = gitlinkCookie

	var attachmentIds []string
	var attachments = make(map[string]interface{})

	for _, artifact := range artifacts {
		attachmentsData, _, err := gitClient.Attachments.PostAttachments(artifact, "")
		if err != nil {
			return nil, nil, err
		}
		if attachmentsData.Status == nil || *attachmentsData.Status == 0 {
			attachmentIds = append(attachmentIds, attachmentsData.Id)

			fileName := filepath.Base(artifact)
			attachments[fileName] = attachmentsData.Url
		} else {
			return nil, nil, errors.New(attachmentsData.Message)
		}
	}

	if gitlinkExportAssetsFileName != "" {

		jsonData, err := json.Marshal(attachments)
		if err != nil {
			log.Println("Error marshal JSON:", err)
			return nil, nil, err
		}

		file, err := os.Create(gitlinkExportAssetsFileName)
		if err != nil {
			log.Printf("Create %s Error:\n%s", gitlinkExportAssetsFileName, err)
			return nil, nil, err
		}
		defer file.Close()

		// 将 JSON 数据写入文件
		_, err = file.Write(jsonData)
		if err != nil {
			log.Printf("Write %s Error:\n%s", gitlinkExportAssetsFileName, err)
			return nil, nil, err
		}
	}

	return attachmentIds, attachments, nil
}

// GitlinkReleases
// 发布
func GitlinkReleases(prerelease bool, releaseName string, releaseBody string, tag string, gitlinkRepository string, gitlinkCookie string, attachmentIds []string, draft bool) error {

	gitClient, err := gitlink.NewClient("")
	if err != nil {
		return err
	}

	gitClient.Cookie = gitlinkCookie

	parts := strings.Split(gitlinkRepository, "/")

	var owner string
	var repo string
	for index, part := range parts {
		if index == 0 {
			owner = part
		} else if index == 1 {
			repo = part
		}
	}

	requestPath := &gitlink.PostReleasesRequestPath{
		Owner: owner,
		Repo:  repo,
	}

	requestBody := &gitlink.PostReleasesRequestBody{
		AttachmentIds: attachmentIds,
		Body:          releaseBody,
		Name:          releaseName,
		TagName:       tag,
		Draft:         draft,
		Prerelease:    prerelease,
	}

	postReleases, _, err := gitClient.Releases.PostReleases(requestPath, requestBody)
	if err != nil {
		return err
	}

	if postReleases.Status == 0 {
		log.Printf("GitLink 发布结果：\n%s", postReleases.Message)
		return nil
	} else {
		return errors.New(fmt.Sprintf("GitLink 发布失败：\n%s", postReleases.Message))
	}
}
