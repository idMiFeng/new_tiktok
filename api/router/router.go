package router

import (
	"api/global"
	"api/internal/social"
	"api/internal/user"
	"api/internal/video"
	"api/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(ErrorMiddleWare())
	// 测试
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	userHandler := user.NewHandlerUser()
	videoHandler := video.NewHandlerUser()
	socialHandler := social.NewHandlerUser()

	baseGroup := r.Group("/douyin")
	{
		// 用户
		baseGroup.POST("/user/register/", userHandler.UserRegister)
		baseGroup.POST("/user/login/", userHandler.UserLogin)
		baseGroup.GET("/user/", utils.JwtToken(), userHandler.UserInfo)
		baseGroup.GET("/feed/", videoHandler.Feed)
		// 视频
		publishGroup := baseGroup.Group("/publish")
		publishGroup.Use(utils.JwtToken())
		{
			publishGroup.POST("/action/", videoHandler.PublishAction)
			publishGroup.GET("/list/", videoHandler.PublishList)
		}
		favoriteGroup := baseGroup.Group("favorite")
		favoriteGroup.Use(utils.JwtToken())
		{
			favoriteGroup.POST("action/", videoHandler.FavoriteAction)
			favoriteGroup.GET("list/", videoHandler.FavoriteList)
		}
		commentGroup := baseGroup.Group("/comment")
		commentGroup.Use(utils.JwtToken())
		{
			commentGroup.POST("/action/", videoHandler.CommentAction)
			commentGroup.GET("/list/", videoHandler.CommentList)
		}

		// 社交
		relationGroup := baseGroup.Group("/relation")
		relationGroup.Use(utils.JwtToken())
		{
			//relationGroup.POST("/action/", socialHandler.FollowAction)
			relationGroup.GET("/follow/list/", socialHandler.GetFollowList)
			relationGroup.GET("/follower/list/", socialHandler.GetFollowerList)
			relationGroup.GET("/friend/list/", socialHandler.GetFriendList)
		}
		messageGroup := baseGroup.Group("/message")
		messageGroup.Use(utils.JwtToken())
		{
			messageGroup.POST("/action/", socialHandler.PostMessage)
			messageGroup.GET("/chat/", socialHandler.GetMessage)
		}
	}

	return r
}

// ErrorMiddleWare 错误处理中间件，捕获panic抛出异常
func ErrorMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r != nil {
				c.JSON(http.StatusOK, gin.H{
					"status_code": global.ErrBadRequest,
					// 打印具体错误
					"status_msg": fmt.Sprintf("%s", r),
				})
				// 中断
				c.Abort()
			}
		}()
		c.Next()
	}
}
