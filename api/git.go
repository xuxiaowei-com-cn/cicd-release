package api

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"os/exec"
	"strings"
)

// RandomString 随机字符串
func RandomString(charset string, length int) string {
	// abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789

	result := ""
	charsetLength := big.NewInt(int64(len(charset)))
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			panic(err)
		}
		result += string(charset[randomIndex.Int64()])
	}
	return result
}

// LowerCaseRandomString 小写随机字符串
func LowerCaseRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz"
	return RandomString(charset, length)
}

// GitVersion 获取 Git 版本号
func GitVersion() (string, error) {
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("无法获取 Git 版本号")
		return "", err
	}

	version := string(output)
	log.Println("Git 版本:", version)
	return version, nil
}

// GitPushTag 推送标签
func GitPushTag(instance string, repository string, username, token, tag string) error {

	origin := LowerCaseRandomString(6)

	instanceUrl, err := url.Parse(instance)
	if err != nil {
		log.Println("Error parsing URL:", err)
		return err
	}

	instanceUrl.User = url.UserPassword(username, token)
	urlUserStr := fmt.Sprintf("%s/%s%s", instanceUrl, repository, ".git")
	urlStr := fmt.Sprintf("%s/%s%s", instance, repository, ".git")

	cmdAddRemote := exec.Command("git", "remote", "add", origin, urlUserStr)
	_, err = cmdAddRemote.Output()
	if err != nil {
		log.Printf("Git 添加临时远端 %s 地址 %s 异常\n", origin, urlStr)
		return err
	}

	log.Printf("Git 推送远端 标签 %s 开始\n", tag)

	cmdPush := exec.Command("git", "push", origin, tag)
	_, err = cmdPush.Output()
	if err != nil {
		log.Printf("Git 推送远端 %s 标签 %s 异常：\n%s", origin, tag, err)
		return err
	}

	log.Printf("Git 推送远端 标签 %s 完成\n", tag)

	cmdRmRemote := exec.Command("git", "remote", "rm", origin)
	_, err = cmdRmRemote.Output()
	if err != nil {
		log.Printf("Git 删除临时远端 %s 地址 %s 异常\n", origin, urlStr)
		return err
	}

	return nil
}

func GitTagSha(tag string) (string, error) {
	cmd := exec.Command("git", "rev-parse", tag)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	sha := strings.ReplaceAll(string(output), "\n", "")

	log.Printf("Git 标签 %s SHA: %s", tag, sha)

	return sha, nil
}

func GitPrintTag(tag string) error {
	cmd := exec.Command("git", "show", tag, "--no-patch")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Git 标签 %s 不存在\n", tag)
		return err
	}

	log.Printf("Git 标签 %s 信息: \n%s", tag, output)
	return nil
}

func GitCreateTag(tag string) error {
	cmd := exec.Command("git", "tag", tag)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("Git 创建标签 %s 异常\n", tag)
		return err
	}

	return GitPrintTag(tag)
}

// AutoCreateTag 自动创建标签
func AutoCreateTag(tag string, autoCreateTag bool) error {
	err := GitPrintTag(tag)
	if err != nil {
		if autoCreateTag {
			log.Printf("开始 创建标签：%s\n", tag)
			err = GitCreateTag(tag)
			if err != nil {
				return err
			}
			log.Printf("完成 创建标签：%s\n", tag)

			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
