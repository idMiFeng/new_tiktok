package social

import (
	"api/internal/res"
	"api/pb/user"
	"strconv"
)

// BuildUsers 构建用户信息
func BuildUsers(userInfos []*user.User) []res.User {

	var userList []res.User

	for i := range userInfos {
		u := &res.User{
			Id:              userInfos[i].Id,
			Name:            userInfos[i].Name,
			FollowCount:     userInfos[i].FollowCount,
			FollowerCount:   userInfos[i].FollowerCount,
			IsFollow:        userInfos[i].IsFollow,
			Avatar:          userInfos[i].Avatar,
			BackgroundImage: userInfos[i].BackgroundImage,
			Signature:       userInfos[i].Signature,
			TotalFavorited:  strconv.FormatInt(userInfos[i].TotalFavorited, 10),
			WorkCount:       userInfos[i].WorkCount,
			FavoriteCount:   userInfos[i].FavoriteCount,
		}

		userList = append(userList, *u)
	}

	return userList
}
