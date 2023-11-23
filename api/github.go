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
	"os"
	"path"
	"strings"
)

type GithubReleasesRequest struct {
	TagName              string `json:"tag_name"`
	TargetCommitish      string `json:"target_commitish"`
	Name                 string `json:"name"`
	Body                 string `json:"body"`
	Draft                bool   `json:"draft"`
	Prerelease           bool   `json:"prerelease"`
	GenerateReleaseNotes bool   `json:"generate_release_notes"`
}

func Github(prerelease bool, context *cli.Context) error {
	var releaseName = context.String(constant.ReleaseName)
	var releaseBody = context.String(constant.ReleaseBody)
	var tag = context.String(constant.Tag)
	var draft = context.Bool(constant.Draft)
	var autoCreateTag = context.Bool(constant.AutoCreateTag)
	var artifacts = context.StringSlice(constant.Artifacts)
	var githubRepository = context.String(constant.GithubRepository)
	var githubUsername = context.String(constant.GithubUsername)
	var githubToken = context.String(constant.GithubToken)

	log.Printf("是否是预发布版本：%v", prerelease)
	log.Printf("发布到 GitHub，路径：%s", githubRepository)

	// 检查发布
	err := GithubGetReleases(githubRepository, githubToken)
	if err != nil {
		return err
	}

	// 自动创建标签
	err = AutoCreateTag(tag, autoCreateTag)
	if err != nil {
		return err
	}

	// 检查标签
	err = GithubGetTag(githubRepository, githubToken, tag)
	if err != nil {
		return err
	}

	// 推送标签
	err = GitPushTag("https://github.com", githubRepository, githubUsername, githubToken, tag)
	if err != nil {
		return err
	}

	// 发布
	err = GithubReleases(prerelease, githubRepository, releaseName, releaseBody, tag, draft, artifacts, githubToken)
	if err != nil {
		return err
	}

	return nil
}

// GithubGetTag
// 检查标签
func GithubGetTag(githubRepository string, githubToken string, tag string) error {

	url := fmt.Sprintf("https://api.github.com/repos/%s/tags/protection", githubRepository)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", githubToken))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

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

	bodyStr := string(body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var tags []Tag
		err := json.Unmarshal(body, &tags)
		if err != nil {
			log.Println("Error unmarshal Tag:", err)
			return err
		}

		for _, t := range tags {
			if tag == t.Mame {
				sha, err := GitTagSha(tag)
				if err != nil {
					return nil
				}

				if t.Commit.Sha == sha {
					return nil
				} else {
					return errors.New(fmt.Sprintf("本地标签 %s（%s） 和 远端 标签 %s（%s） 对应 SHA 不同，请检查！", tag, sha, tag, t.Commit.Sha))
				}
			}
		}

	} else {
		return errors.New(fmt.Sprintf("检查 GitHub 标签异常（%d）：\n%s", resp.StatusCode, bodyStr))
	}

	return nil
}

// GithubGetReleases
// 检查发布
func GithubGetReleases(githubRepository string, githubToken string) error {

	return nil
}

// GithubReleases
// 发布
func GithubReleases(prerelease bool, githubRepository string, releaseName string, releaseBody string, tag string,
	draft bool, artifacts []string, githubToken string) error {

	sha, err := GitTagSha(tag)
	if err != nil {
		return nil
	}

	data := GithubReleasesRequest{
		TagName:              tag,
		TargetCommitish:      sha,
		Name:                 releaseName,
		Body:                 releaseBody,
		Draft:                draft,
		Prerelease:           prerelease,
		GenerateReleaseNotes: false,
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/releases", githubRepository)

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", githubToken))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

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

	bodyStr := string(body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		log.Printf("发布结果：%s\n", bodyStr)

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println("Error unmarshal JSON:", err)
			return err
		}

		// 从map中取出目标值
		idFloat64, ok := data["id"].(float64)
		if !ok {
			log.Fatal("id value not found or not a float64")
		}

		id := int64(idFloat64)

		err = GithubUploadReleaseAssets(githubRepository, artifacts, id, githubToken)
		if err != nil {
			return err
		}

	} else {
		if strings.Contains(bodyStr, "already_exists") {
			return errors.New(fmt.Sprintf("GitHub 已存在发布：\n%s", bodyStr))
		} else {
			return errors.New(fmt.Sprintf("发布 GitHub 异常：\n%s", bodyStr))
		}
	}

	return nil
}

func GithubUploadReleaseAssets(githubRepository string, artifacts []string, id int64, githubToken string) error {

	for _, artifact := range artifacts {

		fileName := path.Base(artifact)

		url := fmt.Sprintf("https://uploads.github.com/repos/%s/releases/%d/assets?name=%s", githubRepository, id, fileName)

		file, err := os.Open(artifact)
		if err != nil {
			log.Printf("Failed to open file: %v", err)
			return err
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("Failed to get file info: %v", err)
			return err
		}

		req, err := http.NewRequest(http.MethodPost, url, file)
		if err != nil {
			log.Printf("Failed to create request: %v", err)
			return err
		}

		req.ContentLength = fileInfo.Size() // 设置正确的Content-Length

		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", githubToken))
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		req.Header.Set("Content-Type", "application/octet-stream")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to send request: %v", err)
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			return err
		}

		bodyStr := string(body)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("上传产物 %s 完成：\n%s\n", artifact, bodyStr)
		} else {
			return errors.New(fmt.Sprintf("上传产物 %s 异常：\n%s", artifact, bodyStr))
		}
	}

	return nil
}
