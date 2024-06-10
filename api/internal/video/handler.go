package video

import (
	"api/global"
	"api/internal/res"
	"api/pb/user"
	"api/pb/video"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"user_service/logger"
)

type HandlerVideo struct {
}

func NewHandlerUser() *HandlerVideo {
	return &HandlerVideo{}
}

// Feed 视频流
func (*HandlerVideo) Feed(c *gin.Context) {
	var feedReq video.FeedRequest

	// 判断是否带有参数
	token := c.Query("token")
	if token == "" {
		feedReq.UserId = -1
	} else {
		userId, _ := c.Get("user_id")
		feedReq.UserId, _ = userId.(int64)
	}

	latestTime := c.Query("latest_time")
	if latestTime == "" || latestTime == "0" {
		feedReq.LatestTime = -1
	} else {
		timePoint, _ := strconv.ParseInt(latestTime, 10, 64)
		feedReq.LatestTime = timePoint
	}
	feedResp, err := global.VideoClient.Feed(context.Background(), &feedReq)
	if err != nil {
		panic(err)
	}
	var userIds []int64
	for _, v := range feedResp.VideoList {
		userIds = append(userIds, v.AuthId)
	}
	userInfos, err := global.UserClient.UserInfo(context.Background(), &user.UserInfoRequest{
		UserIds: userIds,
	})

	list := BuildVideoList(feedResp.VideoList, userInfos.Users)
	r := res.FeedResponse{
		StatusCode: feedResp.StatusCode,
		StatusMsg:  feedResp.StatusMsg,
		NextTime:   feedResp.NextTime,
		VideoList:  list,
	}
	c.JSON(200, r)
}

// PublishAction 发布视频
func (*HandlerVideo) PublishAction(c *gin.Context) {
	var publishActionReq video.PublishActionRequest

	userId, _ := c.Get("user_id")
	publishActionReq.UserId = userId.(int64)

	publishActionReq.Title = c.PostForm("title")

	formFile, _ := c.FormFile("data")
	file, err := formFile.Open()
	if err != nil {
		logger.Log.Error(err)
		res.ErrorResponse(c, 1, err.Error())
		return
	}
	defer file.Close()
	buf, err := io.ReadAll(file) // 将文件读取到字节切片buf中
	if err != nil {
		logger.Log.Error(err)
		res.ErrorResponse(c, 1, err.Error())
		return
	}
	publishActionReq.Data = buf

	videoServiceResp, err := global.VideoClient.PublishAction(context.Background(), &publishActionReq)
	if err != nil {
		logger.Log.Error(err)
		res.ErrorResponse(c, 1, err.Error())
		return
	}

	r := res.PublishActionResponse{
		StatusCode: videoServiceResp.StatusCode,
		StatusMsg:  videoServiceResp.StatusMsg,
	}

	c.JSON(http.StatusOK, r)

}

// PublishList 发布列表
func (*HandlerVideo) PublishList(c *gin.Context) {
	var pulishListReq video.PublishListRequest

	userIdStr := c.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	pulishListReq.UserId = userId
	publishListResp, err := global.VideoClient.PublishList(context.Background(), &pulishListReq)
	if err != nil {
		logger.Log.Error(err)
		res.ErrorResponse(c, 1, err.Error())
		return
	}

	var userIds []int64
	for _, v := range publishListResp.VideoList {
		userIds = append(userIds, v.AuthId)
	}
	userInfos, err := global.UserClient.UserInfo(context.Background(), &user.UserInfoRequest{
		UserIds: userIds,
	})

	// 找到所有的用户信息

	list := BuildVideoList(publishListResp.VideoList, userInfos.Users)

	r := res.VideoListResponse{
		StatusCode: publishListResp.StatusCode,
		StatusMsg:  publishListResp.StatusMsg,
		VideoList:  list,
	}
	c.JSON(200, r)

}

// FavoriteAction 喜欢操作
func (*HandlerVideo) FavoriteAction(c *gin.Context) {
	var favoriteActionReq video.FavoriteActionRequest

	userId, _ := c.Get("user_id")
	favoriteActionReq.UserId, _ = userId.(int64)
	videoId := c.PostForm("video_id")
	favoriteActionReq.VideoId, _ = strconv.ParseInt(videoId, 10, 64)
	actionType := c.PostForm("action_type")
	actionTypeValue, _ := strconv.Atoi(actionType)

	// 异常操作
	if actionTypeValue == 1 || actionTypeValue == 2 {
		favoriteActionReq.ActionType = int64(actionTypeValue)

		videoServiceResp, err := global.VideoClient.FavoriteAction(context.Background(), &favoriteActionReq)
		if err != nil {
			logger.Log.Error(err)
			res.ErrorResponse(c, 1, err.Error())
			return
		}

		r := res.FavoriteActionResponse{
			StatusCode: videoServiceResp.StatusCode,
			StatusMsg:  videoServiceResp.StatusMsg,
		}

		c.JSON(http.StatusOK, r)

	} else {
		res.ErrorResponse(c, 1, "action_type参数错误")
		return
	}
}

// FavoriteList 获取喜欢列表
func (*HandlerVideo) FavoriteList(c *gin.Context) {
	var favoriteListReq video.FavoriteListRequest

	userIdStr := c.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	favoriteListReq.UserId = userId
	favoriteListResp, err := global.VideoClient.FavoriteList(context.Background(), &favoriteListReq)
	if err != nil {
		logger.Log.Error(err)
		res.ErrorResponse(c, 1, err.Error())
		return
	}

	// 找到所有的用户Id
	var userIds []int64
	for _, v := range favoriteListResp.VideoList {
		userIds = append(userIds, v.AuthId)
	}

	// 找到所有的用户信息
	userInfos, err := global.UserClient.UserInfo(context.Background(), &user.UserInfoRequest{
		UserIds: userIds,
	})

	list := BuildVideoList(favoriteListResp.VideoList, userInfos.Users)

	r := res.VideoListResponse{
		StatusCode: favoriteListResp.StatusCode,
		StatusMsg:  favoriteListResp.StatusMsg,
		VideoList:  list,
	}

	c.JSON(http.StatusOK, r)
}

// CommentAction 评论操作
func (*HandlerVideo) CommentAction(c *gin.Context) {
	var commentActionReq video.CommentActionRequest

	userId, _ := c.Get("user_id")
	commentActionReq.UserId, _ = userId.(int64)

	videoId := c.PostForm("video_id")

	commentActionReq.VideoId, _ = strconv.ParseInt(videoId, 10, 64)

	actionType := c.PostForm("action_type")

	actionTypeValue, _ := strconv.Atoi(actionType)
	commentActionReq.ActionType = int64(actionTypeValue)

	// 评论操作
	if commentActionReq.ActionType == 1 {
		commentText := c.PostForm("comment_text")

		commentActionReq.CommentText = commentText
	} else if commentActionReq.ActionType == 2 {
		commentId := c.PostForm("comment_id")
		commentActionReq.CommentId, _ = strconv.ParseInt(commentId, 10, 64)
	} else {
		logger.Log.Error("评论操作类型错误")
		res.ErrorResponse(c, 1, "评论操作类型错误")
		return
	}

	videoServiceResp, err := global.VideoClient.CommentAction(context.Background(), &commentActionReq)
	if err != nil {
		logger.Log.Error(err)
		res.ErrorResponse(c, 1, err.Error())
		return
	}

	if actionTypeValue == 1 {
		// 构建用户信息
		userIds := []int64{userId.(int64)}
		userInfos, _ := global.UserClient.UserInfo(context.Background(), &user.UserInfoRequest{
			UserIds: userIds,
		})

		r := res.CommentActionResponse{
			StatusCode: videoServiceResp.StatusCode,
			StatusMsg:  videoServiceResp.StatusMsg,
			Comment:    BuildComment(videoServiceResp.Comment, userInfos.Users[0]),
		}

		c.JSON(http.StatusOK, r)
	}
	// 如果是删除评论的操作
	if actionTypeValue == 2 {
		r := res.CommentDeleteResponse{
			StatusCode: videoServiceResp.StatusCode,
			StatusMsg:  videoServiceResp.StatusMsg,
		}

		c.JSON(http.StatusOK, r)
	}
}

// CommentList 获取评论列表
func (*HandlerVideo) CommentList(c *gin.Context) {
	var commentListReq video.CommentListRequest

	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)

	commentListReq.VideoId = videoId

	commentListResp, err := global.VideoClient.CommentList(context.Background(), &commentListReq)
	if err != nil {
		logger.Log.Error(err)
		res.ErrorResponse(c, 1, err.Error())
		return
	}

	// 找到所有的用户Id
	var userIds []int64
	for _, comment := range commentListResp.CommentList {
		userIds = append(userIds, comment.UserId)
	}

	userInfos, _ := global.UserClient.UserInfo(context.Background(), &user.UserInfoRequest{
		UserIds: userIds,
	})

	commentList := BuildCommentList(commentListResp.CommentList, userInfos.Users)

	r := res.CommentListResponse{
		StatusCode: commentListResp.StatusCode,
		StatusMsg:  commentListResp.StatusMsg,
		Comments:   commentList,
	}

	c.JSON(http.StatusOK, r)
}
