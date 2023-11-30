FROM golang:1.21-alpine3.18 AS builder

WORKDIR /builder

ADD . .

RUN ls \
    && go version \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download \
    && go test ./... -v \
    && go run main.go \
    && go build -o buildinfo/buildinfo buildinfo/buildinfo.go \
    && go build -ldflags "-s -w -buildid= -X main.BuildDate=$(buildinfo/buildinfo now) -X main.Compiler= -X main.GitCommitBranch=$CI_COMMIT_BRANCH -X main.GitCommitSha=$(buildinfo/buildinfo commitSha) -X main.GitCommitShortSha=$(buildinfo/buildinfo commitShortSha) -X main.GitCommitTag=$(buildinfo/buildinfo commitTag) -X main.GitCommitTimestamp=$(buildinfo/buildinfo commitTimestamp) -X main.GitTreeState=$(buildinfo/buildinfo git-tree-state) -X main.GitVersion=$(buildinfo/buildinfo commitTag) -X main.GoVersion=$(buildinfo/buildinfo goShortVersion) -X main.Major= -X main.Minor= -X main.Revision= -X main.Platform=linux/amd64   -X main.InstanceUrl=$CI_SERVER_URL -X main.CiPipelineId=$CI_PIPELINE_ID -X main.CiJobId=$CI_JOB_ID" -trimpath -o cicd-release . \
    && ./cicd-release --version

FROM alpine:3.18.3

MAINTAINER 徐晓伟 xuxiaowei@xuxiaowei.com.cn

COPY --from=builder /builder/cicd-release /usr/bin/cicd-release
