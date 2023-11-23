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
	"os"
	"path"
)

type Link struct {
	Name            string `json:"name"`
	Url             string `json:"url"`
	DirectAssetPath string `json:"direct_asset_path,omitempty"`
	LinkType        string `json:"link_type,omitempty"`
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
	var packageName = context.String(constant.PackageName)
	var autoCreateTag = context.Bool(constant.AutoCreateTag)
	var milestones = context.StringSlice(constant.Milestones)
	var artifacts = context.StringSlice(constant.Artifacts)
	var gitlabInstance = context.String(constant.GitlabInstance)
	var gitlabApi = context.String(constant.GitlabApi)
	var gitlabRepository = context.String(constant.GitlabRepository)
	var gitlabToken = context.String(constant.GitlabToken)
	var gitlabExportAssetsFileName = context.String(constant.GitlabExportAssetsFileName)

	log.Printf("是否是预发布版本：%v", prerelease)
	log.Printf("发布到 GitLab，实例：%s，路径：%s", gitlabInstance, gitlabRepository)

	baseUrl, err := url.Parse(gitlabInstance)
	if err != nil {
		log.Println("Gitlab 实例配置错误，无法转为 URL")
		panic(err)
	}

	gitlabRepositoryEscape := url.PathEscape(gitlabRepository)

	getReleasesUrl := fmt.Sprintf("%s/%s/projects/%s/releases/%s", baseUrl, gitlabApi, gitlabRepositoryEscape, tag)
	err = GitlabGetReleases(getReleasesUrl, gitlabToken)
	if err != nil {
		return err
	}

	err = AutoCreateTag(tag, autoCreateTag)
	if err != nil {
		return err
	}

	getTagUrl := fmt.Sprintf("%s/%s/projects/%s/repository/tags/%s", baseUrl, gitlabApi, gitlabRepositoryEscape, tag)
	err = GitlabGetTag(getTagUrl, gitlabToken, tag)
	if err != nil {
		return err
	}

	err = GitPushTag(gitlabInstance, gitlabRepository, gitlabToken, tag)
	if err != nil {
		return err
	}

	genericPackagesPrefixUrl := fmt.Sprintf("%s/%s/projects/%s/packages/generic/%s/%s", baseUrl, gitlabApi, gitlabRepositoryEscape, packageName, tag)
	genericPackages, err := GitlabGenericPackages(genericPackagesPrefixUrl, artifacts, gitlabToken, gitlabInstance, gitlabRepository, gitlabExportAssetsFileName)
	if err != nil {
		return err
	}

	err = GitlabReleases(releaseName, releaseBody, tag, milestones,
		baseUrl, gitlabApi, gitlabRepositoryEscape, gitlabToken, genericPackages)
	if err != nil {
		return err
	}

	return nil
}

func GitlabGetTag(getTagUrl string, gitlabToken string, tag string) error {

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

	bodyStr := string(body)

	if resp.StatusCode == 404 {
		return nil
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println("Error unmarshal JSON:", err)
			return err
		}

		// 从map中取出目标值
		targetValue, ok := data["target"].(string)
		if !ok {
			log.Fatal("Target value not found or not a string")
		}

		sha, err := GitTagSha(tag)
		if err != nil {
			return nil
		}

		if sha == targetValue {
			return nil
		} else {
			return errors.New(fmt.Sprintf("本地标签 %s（%s） 和 远端 标签 %s（%s） 对应 SHA 不同，请检查！", tag, sha, tag, targetValue))
		}
	} else {
		return errors.New(fmt.Sprintf("检查远端标签失败：%s", bodyStr))
	}
}

func GitlabGetReleases(getReleasesUrl string, gitlabToken string) error {

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

	bodyStr := string(body)

	if resp.StatusCode == 404 {
		return nil
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return errors.New(fmt.Sprintf("已存在此发布：\n%s", bodyStr))
	} else {
		return errors.New(fmt.Sprintf("检查发布失败：\n%s", bodyStr))
	}
}

func GitlabGenericPackages(genericPackagesPrefixUrl string, artifacts []string, gitlabToken string,
	gitlabInstance string, gitlabRepository string, gitlabExportAssetsFileName string) (map[string]interface{}, error) {

	if artifacts == nil {
		log.Println("未设置上传的产物")
		return nil, nil
	}

	log.Println("开始 上传产物")

	packageFilePrefixUrl := fmt.Sprintf("%s/%s", gitlabInstance, gitlabRepository)

	result := make(map[string]interface{})

	for _, artifact := range artifacts {

		fileName := path.Base(artifact)
		fmt.Println(fileName)

		genericPackagesUrl := fmt.Sprintf("%s/%s?select=package_file", genericPackagesPrefixUrl, fileName)
		log.Printf("上传产物 %s 的 URL %s\n", artifact, genericPackagesUrl)

		file, err := os.Open(artifact)
		if err != nil {
			fmt.Printf("Failed to open file: %s\n", err)
			return nil, err
		}
		defer file.Close()

		req, err := http.NewRequest(http.MethodPut, genericPackagesUrl, file)
		if err != nil {
			fmt.Printf("Failed to create request: %s\n", err.Error())
			return nil, err
		}
		req.Header.Set("PRIVATE-TOKEN", gitlabToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Request failed: %s\n", err.Error())
			return nil, err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println("Error closing response body:", err)

			}
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response body: %s\n", err.Error())
			return nil, err
		}

		bodyStr := string(body)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var data map[string]interface{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				log.Println("Error unmarshal JSON:", err)
				return nil, err
			}

			// 从map中取出目标值
			idFloat64, ok := data["id"].(float64)
			if !ok {
				log.Fatal("Target value not found or not a float64")
			}

			id := int64(idFloat64)

			packageFileUrl := fmt.Sprintf("%s/-/package_files/%d/download", packageFilePrefixUrl, id)

			log.Printf("上传产物 %s 的 下载地址 %s\n", artifact, packageFileUrl)
			result[fileName] = packageFileUrl

		} else {
			return nil, errors.New(fmt.Sprintf("上传产物异常：\n%s", bodyStr))
		}
	}

	log.Println("完成 上传产物")

	if gitlabExportAssetsFileName != "" {
		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Println("Error marshal JSON:", err)
			return nil, err
		}

		file, err := os.Create(gitlabExportAssetsFileName)
		if err != nil {
			log.Printf("Create %s Error:\n%s", gitlabExportAssetsFileName, err)
			return nil, err
		}
		defer file.Close()

		// 将 JSON 数据写入文件
		_, err = file.Write(jsonData)
		if err != nil {
			log.Printf("Write %s Error:\n%s", gitlabExportAssetsFileName, err)
			return nil, err
		}
	}

	return result, nil
}

func GitlabReleases(releaseName string, releaseBody string, tag string, milestones []string,
	baseUrl *url.URL, gitlabApi string, gitlabRepositoryEscape string, gitlabToken string,
	genericPackages map[string]interface{}) error {

	data := Data{
		Name:        releaseName,
		TagName:     tag,
		Description: releaseBody,
		Milestones:  milestones,
	}

	assets := Assets{}
	if genericPackages != nil {
		for key, value := range genericPackages {
			link := Link{
				Name: key,
				Url:  value.(string),
			}
			assets.Links = append(assets.Links, link)
		}
	}
	data.Assets = assets

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

	return nil
}
