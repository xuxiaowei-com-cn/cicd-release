# 说明

此文件夹下的配置用于 [建木Hub](https://www.jianmuhub.com) 环境下构建 Docker 镜像（分阶段构建）

## 建木Hub 构建命令

- 在项目根目录下执行

```shell
# --no-cache：不使用缓存
# --progress plain：显示构建日志
docker build . -f ./builder/jianmuhub/alpine/Dockerfile -t docker.jianmuhub.com/xuxiaowei/cicd-release:alpine --no-cache --progress plain
docker build . -f ./builder/jianmuhub/debian/Dockerfile -t docker.jianmuhub.com/xuxiaowei/cicd-release:debian --no-cache --progress plain
```

## 建木Hub 镜像

- https://res.jianmuhub.com/image/xuxiaowei/cicd-release
