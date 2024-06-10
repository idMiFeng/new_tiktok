package global

import (
	"api/pb/social"
	"api/pb/user"
	"api/pb/video"
)

var (
	UserClient   user.UserServiceClient
	VideoClient  video.VideoServiceClient
	SocialClient social.SocialServiceClient
)
