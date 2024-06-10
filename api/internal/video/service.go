package video

import (
	"api/internal/res"
	"api/pb/user"
	"api/pb/video"
	"strconv"
)

// BuildVideoList 构建视频列表
func BuildVideoList(videos []*video.Video, userInfos []*user.User) []res.Video {

	var videoList []res.Video

	for i, v := range videos {
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
		videoList = append(videoList, res.Video{
			Id:            v.Id,
			Author:        *u,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
		})
	}

	return videoList
}

// BuildComment 构建单个评论
func BuildComment(comment *video.Comment, userInfo *user.User) res.Comment {
	u := &res.User{
		Id:              userInfo.Id,
		Name:            userInfo.Name,
		FollowCount:     userInfo.FollowCount,
		FollowerCount:   userInfo.FollowerCount,
		IsFollow:        userInfo.IsFollow,
		Avatar:          userInfo.Avatar,
		BackgroundImage: userInfo.BackgroundImage,
		Signature:       userInfo.Signature,
		TotalFavorited:  strconv.FormatInt(userInfo.TotalFavorited, 10),
		WorkCount:       userInfo.WorkCount,
		FavoriteCount:   userInfo.FavoriteCount,
	}
	return res.Comment{
		Id:         comment.Id,
		User:       *u,
		Content:    comment.Content,
		CreateDate: comment.CreateDate,
	}
}

// BuildCommentList 构建评论列表
func BuildCommentList(comments []*video.Comment, userInfos []*user.User) []res.Comment {
	var commentList []res.Comment

	for i, comment := range comments {
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
		commentList = append(commentList, res.Comment{
			Id:         comment.Id,
			User:       *u,
			Content:    comment.Content,
			CreateDate: comment.CreateDate,
		})
	}

	return commentList
}
