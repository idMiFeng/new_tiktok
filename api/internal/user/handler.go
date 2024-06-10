package user

import (
	"api/config"
	"api/global"
	"api/internal/res"
	"api/pb/user"
	"api/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
	"user_service/logger"
)

type HandlerUser struct {
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{}
}
func BindErrorResponse(c *gin.Context, err error) {
	if ves, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ves {
			c.JSON(http.StatusOK, user.UserResponse{
				StatusCode: int64(global.ErrBadRequest),
				StatusMsg:  fe.Translate(utils.GetTrans()),
			})
			log.Println(fe, "|", fe.Translate(utils.GetTrans()))
			return
		}
	}
	ErrorResponse(c, "参数绑定失败", err)
}
func ErrorResponse(c *gin.Context, msg string, err error) {
	c.JSON(http.StatusOK, user.UserResponse{
		StatusCode: int64(global.ErrBadRequest),
		StatusMsg:  msg,
	})
	log.Printf("%s: %v\n", msg, err)
}

// UserRegister 用户注册
func (h *HandlerUser) UserRegister(c *gin.Context) {
	// 处理注册逻辑
	var userReq user.UserRequest
	if err := c.ShouldBind(&userReq); err != nil {
		BindErrorResponse(c, err)
		return
	}
	userServiceClient := global.UserClient
	userResp, err := userServiceClient.UserRegister(context.Background(), &userReq)
	if err != nil {
		logger.Log.Error("UserService--error:", err)
	}
	j := utils.NewJWT(config.Conf.JwtSecret)
	claims := utils.MyClaims{
		UserId: strconv.FormatInt(userResp.UserId, 10),
	}
	token, err := j.CreateToken(claims)
	r := res.UserResponse{
		StatusCode: userResp.StatusCode,
		StatusMsg:  userResp.StatusMsg,
		UserId:     userResp.UserId,
		Token:      token,
	}
	c.JSON(http.StatusOK, r)
}

// UserLogin 用户登录
func (h *HandlerUser) UserLogin(c *gin.Context) {
	var userReq user.UserRequest
	if err := c.ShouldBind(&userReq); err != nil {
		BindErrorResponse(c, err)
		return
	}
	userServiceClient := global.UserClient
	userResp, err := userServiceClient.UserLogin(context.Background(), &userReq)
	if err != nil {
		logger.Log.Error("UserService--Login error:", err)
	}
	j := utils.NewJWT(config.Conf.JwtSecret)
	claims := utils.MyClaims{
		UserId: strconv.FormatInt(userResp.UserId, 10),
	}
	token, err := j.CreateToken(claims)
	r := res.UserResponse{
		StatusCode: userResp.StatusCode,
		StatusMsg:  userResp.StatusMsg,
		UserId:     userResp.UserId,
		Token:      token,
	}
	c.JSON(http.StatusOK, r)
}

// UserInfo 用户信息列表
func (h *HandlerUser) UserInfo(c *gin.Context) {
	var userIds []int64

	// jwt中间件会解析token，然后把user_id放入context中，所以用两种方式都可以获取到user_id
	Id, _ := c.Get("user_id")
	userIdStr := Id.(string)
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	userIds = append(userIds, userId)
	userServiceClient := global.UserClient
	userResp, err := userServiceClient.UserInfo(context.Background(), &user.UserInfoRequest{UserIds: userIds})
	if err != nil {
		logger.Log.Error("UserService--UserInfo error:", err)
	}
	var userinfo res.User
	userinfo.Id = userResp.Users[0].Id
	userinfo.Name = userResp.Users[0].Name
	userinfo.FollowCount = userResp.Users[0].FollowCount
	userinfo.FollowerCount = userResp.Users[0].FollowerCount
	userinfo.IsFollow = userResp.Users[0].IsFollow
	userinfo.Avatar = userResp.Users[0].Avatar
	userinfo.BackgroundImage = userResp.Users[0].BackgroundImage
	userinfo.Signature = userResp.Users[0].Signature
	userinfo.TotalFavorited = strconv.FormatInt(userResp.Users[0].TotalFavorited, 10)
	userinfo.WorkCount = userResp.Users[0].WorkCount
	userinfo.FavoriteCount = userResp.Users[0].FavoriteCount

	r := res.UserInfoResponse{
		StatusCode: userResp.StatusCode,
		StatusMsg:  userResp.StatusMsg,
		User:       userinfo,
	}
	c.JSON(http.StatusOK, r)
}
