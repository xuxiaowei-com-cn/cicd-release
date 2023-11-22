package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xuxiaowei-com-cn/cicd-release/constant"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Link struct {
	Name            string `json:"name"`
	Url             string `json:"url"`
	DirectAssetPath string `json:"direct_asset_path"`
	LinkType        string `json:"link_type"`
}

type Assets struct {
	Links []Link `json:"links"`
}

type Data struct {
	Name        string   `json:"name"`
	TagName     string   `json:"tag_name"`
	Description string   `json:"description"`
	Milestones  []string `json:"milestones"`
	Assets      Assets   `json:"assets"`
}

func Gitlab(prerelease bool, context *cli.Context) error {
	var releaseName = context.String(constant.ReleaseName)
	var releaseBody = context.String(constant.ReleaseBody)
	var tag = context.String(constant.Tag)
	var autoCreateTag = context.Bool(constant.AutoCreateTag)
	var milestones = context.StringSlice(constant.Milestones)
	var artifacts = context.StringSlice(constant.Artifacts)
	var gitlabInstance = context.String(constant.GitlabInstance)
	var gitlabApi = context.String(constant.GitlabApi)
	var gitlabRepository = context.String(constant.GitlabRepository)
	var gitlabToken = context.String(constant.GitlabToken)

	log.Printf("是否是预发布版本：%v", prerelease)
	log.Printf("发布到 GitLab，实例：%s", gitlabInstance)

	baseUrl, err := url.Parse(gitlabInstance)
	if err != nil {
		log.Println("Gitlab 实例配置错误，无法转为 URL")
		panic(err)
	}

	gitlabRepositoryEscape := url.PathEscape(gitlabRepository)

	getReleasesUrl := fmt.Sprintf("%s/%s/%s/%s/releases/%s", baseUrl, gitlabApi, "projects", gitlabRepositoryEscape, tag)

	err = GetReleases(getReleasesUrl, gitlabToken)
	if err != nil {
		return err
	}

	err = AutoCreateTag(tag, autoCreateTag)
	if err != nil {
		return err
	}

	err = GitPushTag(gitlabInstance, gitlabRepository, gitlabToken, tag)
	if err != nil {
		return err
	}

	releasesUrl := fmt.Sprintf("%s/%s/%s/%s/releases", baseUrl, gitlabApi, "projects", gitlabRepositoryEscape)

	data := Data{
		Name:        releaseName,
		TagName:     tag,
		Description: releaseBody,
		Milestones:  milestones,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return err
	}

	req, err := http.NewRequest("POST", releasesUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PRIVATE-TOKEN", gitlabToken)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}(resp.Body)

	log.Println("Response status:", resp.Status)

	log.Printf("artifacts：%s", artifacts)

	return nil
}

func GetReleases(getReleasesUrl string, gitlabToken string) error {

	client := &http.Client{}
	req, err := http.NewRequest("GET", getReleasesUrl, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("PRIVATE-TOKEN", gitlabToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return err
	}

	return errors.New(fmt.Sprintf("已存在此发布：\n%s", string(body)))
}
