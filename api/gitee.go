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
)

type Tag struct {
	Mame    string `json:"name"`
	Message string `json:"message"`
	Commit  Commit `json:"commit"`
}

type Commit struct {
	Sha  string `json:"sha"`
	Date string `json:"date"`
}

type GiteeReleasesRequest struct {
	AccessToken     string `json:"access_token"`
	TagName         string `json:"tag_name"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	Prerelease      bool   `json:"prerelease"`
	TargetCommitish string `json:"target_commitish"`
}

func Gitee(prerelease bool, context *cli.Context) error {
	var releaseName = context.String(constant.ReleaseName)
	var releaseBody = context.String(constant.ReleaseBody)
	var tag = context.String(constant.Tag)
	var autoCreateTag = context.Bool(constant.AutoCreateTag)
	var giteeRepository = context.String(constant.GiteeRepository)
	var giteeUsername = context.String(constant.GiteeUsername)
	var giteeToken = context.String(constant.GiteeToken)
	var gitlabExportAssetsFileName = context.String(constant.GitlabExportAssetsFileName)

	log.Printf("是否是预发布版本：%v", prerelease)
	log.Printf("发布到 Gitee，路径：%s", giteeRepository)

	// 检查发布
	err := GiteeGetReleases(giteeRepository, giteeToken)
	if err != nil {
		return err
	}

	// 自动创建标签
	err = AutoCreateTag(tag, autoCreateTag)
	if err != nil {
		return err
	}

	// 检查标签
	err = GiteeGetTag(giteeRepository, giteeToken, tag)
	if err != nil {
		return err
	}

	// 推送标签
	err = GitPushTag("https://gitee.com", giteeRepository, giteeUsername, giteeToken, tag)
	if err != nil {
		return err
	}

	// 发布
	err = GiteeReleases(prerelease, giteeRepository, releaseName, releaseBody, tag, giteeToken, gitlabExportAssetsFileName)
	if err != nil {
		return err
	}

	return nil
}

// GiteeGetTag
// 检查标签
func GiteeGetTag(giteeRepository string, giteeToken string, tag string) error {

	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/tags?access_token=%s", giteeRepository, giteeToken)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Failed to create request: %s", err)
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %s", err)
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
		log.Printf("Failed to read response body: %s\n", err)
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
		return errors.New(fmt.Sprintf("检查 Gitee 标签异常：\n%s", bodyStr))
	}

	return nil
}

// GiteeGetReleases
// 检查发布
func GiteeGetReleases(giteeRepository, giteeToken string) error {

	return nil
}

// GiteeReleases
// 发布
func GiteeReleases(prerelease bool, giteeRepository string, releaseName string, releaseBody string, tag string,
	giteeToken string, gitlabExportAssetsFileName string) error {

	if gitlabExportAssetsFileName != "" {
		jsonData, err := os.ReadFile(gitlabExportAssetsFileName)
		if err != nil {
			log.Printf("ReadFile %s Error:\n%s", gitlabExportAssetsFileName, err)
			return err
		}

		readResult := make(map[string]interface{})

		err = json.Unmarshal(jsonData, &readResult)
		if err != nil {
			log.Printf("Unmarshal %s Error:\n%s", gitlabExportAssetsFileName, err)
			return err
		}

		releaseBody += "\n"
		releaseBody += "***\n"

		for key, value := range readResult {
			releaseBody += fmt.Sprintf("\n- [%s](%s)\n", key, value)
		}
	}

	sha, err := GitTagSha(tag)
	if err != nil {
		return nil
	}

	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/releases", giteeRepository)

	data := GiteeReleasesRequest{
		AccessToken:     giteeToken,
		TagName:         tag,
		Name:            releaseName,
		Body:            releaseBody,
		Prerelease:      prerelease,
		TargetCommitish: sha,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
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

	log.Println("Response status:", resp.Status)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("Gitee 发布结果：\n%s", bodyStr)
		return nil
	} else {
		return errors.New(fmt.Sprintf("Gitee 发布失败：\n%s", bodyStr))
	}
}
