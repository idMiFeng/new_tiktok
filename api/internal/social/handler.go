package social

import (
	"api/global"
	"api/internal/res"
	"api/logger"
	"api/pb/social"
	"api/pb/user"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type HandlerSocial struct {
}

func NewHandlerUser() *HandlerSocial {
	return &HandlerSocial{}
}

// GetFollowList 关注列表
func (*HandlerSocial) GetFollowList(c *gin.Context) {
	var followList social.FollowListRequest
	userId := c.Query("user_id")
	followList.UserId, _ = strconv.ParseInt(userId, 10, 64)

	socialResp, err := global.SocialClient.GetFollowList(context.Background(), &followList)
	if err != nil {
		logger.Log.Error(err)
		res.ErrorResponse(c, 1, err.Error())
		return
	}

	userInfos, err := global.UserClient.UserInfo(context.Background(), &user.UserInfoRequest{
		UserIds: socialResp.UserId,
	})
	r := res.FollowListResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
		UserList:   BuildUsers(userInfos.Users),
	}
	c.JSON(http.StatusOK, r)
}

// GetFollowerList 粉丝列表
func (*HandlerSocial) GetFollowerList(c *gin.Context) {
	var followerList social.FollowListRequest
	userId := c.Query("user_id")
	followerList.UserId, _ = strconv.ParseInt(userId, 10, 64)

	socialResp, err := global.SocialClient.GetFollowerList(context.Background(), &followerList)
	if err != nil {
		logger.Log.Error(err)
		res.ErrorResponse(c, 1, err.Error())
		return
	}
	userInfos, err := global.UserClient.UserInfo(context.Background(), &user.UserInfoRequest{
		UserIds: socialResp.UserId,
	})
	r := res.FollowListResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
		UserList:   BuildUsers(userInfos.Users),
	}
	c.JSON(http.StatusOK, r)
}

// GetFriendList 好友列表
func (*HandlerSocial) GetFriendList(c *gin.Context) {
	var friendList social.FollowListRequest
	userId := c.Query("user_id")
	friendList.UserId, _ = strconv.ParseInt(userId, 10, 64)

	socialResp, err := global.SocialClient.GetFriendList(context.Background(), &friendList)
	if err != nil {
		logger.Log.Error(err)
		res.ErrorResponse(c, 1, err.Error())
		return
	}
	userInfos, err := global.UserClient.UserInfo(context.Background(), &user.UserInfoRequest{
		UserIds: socialResp.UserId,
	})
	r := res.FollowListResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
		UserList:   BuildUsers(userInfos.Users),
	}
	c.JSON(http.StatusOK, r)
}

// PostMessage 发送消息
func (*HandlerSocial) PostMessage(c *gin.Context) {
	var postMessage social.PostMessageRequest
	userId, _ := c.Get("user_id")
	postMessage.UserId, _ = userId.(int64)
	toUserId := c.Query("to_user_id")
	postMessage.ToUserId, _ = strconv.ParseInt(toUserId, 10, 64)
	actionType := c.Query("action_type")
	actionTypeInt64, _ := strconv.ParseInt(actionType, 10, 32)

	postMessage.ActionType = int32(actionTypeInt64)
	content := c.Query("content")
	postMessage.Content = content

	socialResp, err := global.SocialClient.PostMessage(context.Background(), &postMessage)
	if err != nil {
		logger.Log.Error(err)
		res.ErrorResponse(c, 1, err.Error())
		return
	}

	r := res.PostMessageResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
	}
	c.JSON(http.StatusOK, r)
}

// GetMessage 获取消息列表
func (*HandlerSocial) GetMessage(c *gin.Context) {
	var getMessage social.GetMessageRequest
	userId, _ := c.Get("user_id")
	getMessage.UserId, _ = userId.(int64)
	toUserId := c.Query("to_user_id")
	getMessage.ToUserId, _ = strconv.ParseInt(toUserId, 10, 64)
	PreMsgTime := c.Query("pre_msg_time")
	getMessage.PreMsgTime, _ = strconv.ParseInt(PreMsgTime, 10, 64)

	socialResp, err := global.SocialClient.GetMessage(context.Background(), &getMessage)

	if err != nil {
		logger.Log.Error(err)
		res.ErrorResponse(c, 1, err.Error())
		return
	}

	r := new(res.GetMessageResponse)
	r.StatusCode = socialResp.StatusCode
	r.StatusMsg = socialResp.StatusMsg
	for _, message := range socialResp.MessageList {
		messageResp := res.Message{
			Id:         message.Id,
			ToUserId:   message.ToUserId,
			FromUserID: message.UserId,
			Content:    message.Content,
			CreateTime: message.CreatedAt,
		}
		r.MessageList = append(r.MessageList, messageResp)
	}

	c.JSON(http.StatusOK, r)
}
