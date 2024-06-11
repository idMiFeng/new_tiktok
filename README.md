# new_tiktok

## 项目背景

本项目是第六届字节跳动青训营的大项目作业，旨在模拟实现一个类似抖音的短视频应用。项目采用微服务架构，共分为四层，主要使用了Gin和gRPC框架。

## 项目结构

该项目由以下几部分组成：

- **API 层**：对外暴露 HTTP 接口，并调用各个微服务。
- **User Service**：处理用户信息相关的操作。
- **Social Service**：处理用户之间的社交操作，如关注、发消息等。
- **Video Service**：处理视频相关的操作，如刷视频、投稿视频、点赞视频、评论视频等。

## 功能实现

1. **用户信息管理**
   - 用户注册
   - 用户登录
   - 获取用户信息
   - 修改用户信息

2. **视频管理**
   - 浏览视频流
   - 上传视频
   - 点赞视频
   - 评论视频

3. **社交管理**
   - 关注其他用户
   - 给好友发送消息

## 使用说明

### 环境配置
- 需要修改组件的配置，在每个目录下的config目录修改
- 确保已安装 Go 语言环境（版本 1.16 及以上）
- 安装依赖并启动各个微服务：
  ```sh
  cd user_service
  go mod tidy
  go run main.go
  ```
  ```sh
  cd video_service
  go mod tidy
  go run main.go
  ```
  ```sh
  cd social_service
  go mod tidy
  go run main.go
  ```
- 启动 API 层
  ```sh
  cd api
  go mod tidy
  go run main.go
  ```
感谢你的使用！
