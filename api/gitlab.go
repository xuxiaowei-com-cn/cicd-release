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
	"strings"
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
	log.Printf("发布到 GitLab，实例：%s，路径：%s", gitlabInstance, gitlabRepository)

	baseUrl, err := url.Parse(gitlabInstance)
	if err != nil {
		log.Println("Gitlab 实例配置错误，无法转为 URL")
		panic(err)
	}

	gitlabRepositoryEscape := url.PathEscape(gitlabRepository)

	getReleasesUrl := fmt.Sprintf("%s/%s/projects/%s/releases/%s", baseUrl, gitlabApi, gitlabRepositoryEscape, tag)
	err = GetReleases(getReleasesUrl, gitlabToken)
	if err != nil {
		return err
	}

	err = AutoCreateTag(tag, autoCreateTag)
	if err != nil {
		return err
	}

	getTagUrl := fmt.Sprintf("%s/%s/projects/%s/repository/tags/%s", baseUrl, gitlabApi, gitlabRepositoryEscape, tag)
	err = GetTag(getTagUrl, gitlabToken, tag)
	if err != nil {
		return err
	}

	err = GitPushTag(gitlabInstance, gitlabRepository, gitlabToken, tag)
	if err != nil {
		return err
	}

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

	releasesUrl := fmt.Sprintf("%s/%s/projects/%s/releases", baseUrl, gitlabApi, gitlabRepositoryEscape)
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

func GetTag(getTagUrl string, gitlabToken string, tag string) error {

	client := &http.Client{}
	req, err := http.NewRequest("GET", getTagUrl, nil)
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

	if "{\"message\":\"404 Tag Not Found\"}" == string(body) {
		return nil
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Unmarshal Error:", err)
		return err
	}

	// 从map中取出目标值
	targetValue, ok := data["target"].(string)
	if !ok {
		log.Println("Target value not found or not a string")
	}

	sha, err := GitTagSha(tag)
	if err != nil {
		return nil
	}

	if sha != targetValue {
		return errors.New(fmt.Sprintf("本地标签 %s（%s） 和 远端 标签 %s（%s） 对应 SHA 不同，请检查！", tag, strings.ReplaceAll(sha, "\n", ""), tag, targetValue))
	}

	return errors.New(fmt.Sprintf("已存在此标签：\n%s", string(body)))
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

	if "{\"message\":\"404 Not Found\"}" == string(body) {
		return nil
	}

	return errors.New(fmt.Sprintf("已存在此发布：\n%s", string(body)))
}
