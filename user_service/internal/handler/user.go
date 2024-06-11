package handler

import (
	"api/pb/user"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"user_service/global"
	"user_service/internal/dao"
	"user_service/internal/model"
)

type UserService struct {
	user.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (*UserService) UserRegister(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	resp = new(user.UserResponse)
	var userinfo model.User
	// 检查用户是否已经存在
	if exist, err := model.GetUserInstance().CheckUserExist(req.Username); !exist {
		resp.StatusCode = global.ErrCodeUsernameAlreadyExist
		resp.StatusMsg = global.GetErrorMessage(global.ErrCodeUsernameAlreadyExist)
		return resp, err
	}
	userinfo.Username = req.Username
	userinfo.Password = req.Password
	err = model.GetUserInstance().CreateUserInfo(&userinfo)
	if err != nil {
		resp.StatusCode = global.ErrCodeFailedToCreateUser
		resp.StatusMsg = global.GetErrorMessage(global.ErrCodeFailedToCreateUser)
		return resp, err
	}
	// 查询出ID
	userName, err := model.GetUserInstance().GetUserInfoByUsername(userinfo.Username)
	if err != nil {
		resp.StatusCode = global.ErrCodeUserNotFound
		resp.StatusMsg = global.GetErrorMessage(global.ErrCodeUserNotFound)
		return resp, err
	}

	resp.StatusCode = global.SuccessCode
	resp.StatusMsg = global.GetErrorMessage(global.SuccessCode)
	resp.UserId = userName.ID
	return resp, nil
}

func (*UserService) UserLogin(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	resp = new(user.UserResponse)
	// 检查用户是否已经存在
	if exist, err := model.GetUserInstance().CheckUserExist(req.Username); !exist {
		resp.StatusCode = global.ErrCodeUsernameAlreadyExist
		resp.StatusMsg = global.GetErrorMessage(global.ErrCodeUsernameAlreadyExist)
		return resp, err
	}

	userinfo, err := model.GetUserInstance().GetUserInfoByUsername(req.Username)
	// 检查密码是否正确
	if ok := model.GetUserInstance().CheckPassWord(req.Password, userinfo.Password); !ok {
		resp.StatusCode = global.ErrCodeInvalidUsernameOrPass
		resp.StatusMsg = global.GetErrorMessage(global.ErrCodeInvalidUsernameOrPass)
		return resp, err
	}
	resp.StatusCode = global.SuccessCode
	resp.StatusMsg = global.GetErrorMessage(global.SuccessCode)
	resp.UserId = userinfo.ID
	return resp, nil
}

// UserInfo 获取用户信息
func (*UserService) UserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	resp = new(user.UserInfoResponse)
	userIds := req.UserIds
	for _, userId := range userIds {
		// 查看缓存是否存在 需要的信息
		var userinfo *model.User
		var usercount *model.UserCount
		key := fmt.Sprintf("%s:%s:%s", "user", "info", strconv.FormatInt(userId, 10))
		exists, err := global.Redis.Exists(dao.Ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("userinfo缓存错误：%v", err)
		}
		key2 := fmt.Sprintf("%s:%s:%s", "user", "count", strconv.FormatInt(userId, 10))
		exists2, err := global.Redis.Exists(dao.Ctx, key2).Result()
		if err != nil {
			return nil, fmt.Errorf("usercount缓存错误：%v", err)
		}

		if exists > 0 {
			// 查询缓存
			userString, err := global.Redis.Get(dao.Ctx, key).Result()
			if err != nil {
				return nil, fmt.Errorf("userinfo缓存错误：%v", err)
			}
			err = json.Unmarshal([]byte(userString), &userinfo)
			if err != nil {
				return nil, err
			}
		} else {
			// 查询数据库
			userinfo, err = model.GetUserInstance().GetUserInfoByUserID(userId)
			if err != nil {
				resp.StatusCode = global.ErrCodeFailedToGetUserInfo
				resp.StatusMsg = global.GetErrorMessage(global.ErrCodeFailedToGetUserInfo)
				return resp, err
			}
			// 将查询结果放入缓存中
			userJson, _ := json.Marshal(&userinfo)
			err = global.Redis.Set(dao.Ctx, key, userJson, 12*time.Hour).Err()
			if err != nil {
				return nil, fmt.Errorf("缓存错误：%v", err)
			}

		}

		if exists2 > 0 {
			// 查询缓存
			userString, err := global.Redis.Get(dao.Ctx, key2).Result()
			if err != nil {
				return nil, fmt.Errorf("usercount缓存错误：%v", err)
			}
			err = json.Unmarshal([]byte(userString), &usercount)
			if err != nil {
				return nil, err
			}
		} else {
			// 查询数据库
			usercount, err = model.GetUserInstance().GetUserCountByUserID(userId)
			if err != nil {
				resp.StatusCode = global.ErrCodeFailedToGetUserInfo
				resp.StatusMsg = global.GetErrorMessage(global.ErrCodeFailedToGetUserInfo)
				return resp, err
			}
			// 将查询结果放入缓存中
			userJson, _ := json.Marshal(&usercount)
			err = global.Redis.Set(dao.Ctx, key2, userJson, 12*time.Hour).Err()
			if err != nil {
				return nil, fmt.Errorf("缓存错误：%v", err)
			}

		}
		resp.Users = append(resp.Users, BuildUser(userinfo, usercount))
	}
	return resp, nil
}

func (*UserService) UpdateUserFav(ctx context.Context, req *user.UpdateUserFavRequest) (resp *user.UpdateUserFavResponse, err error) {
	resp = new(user.UpdateUserFavResponse)
	var usercount *model.UserCount
	key := fmt.Sprintf("%s:%s:%s", "user", "count", strconv.FormatInt(req.UserId, 10))

	// 检查 key 是否存在
	exists, err := global.Redis.Exists(ctx, key).Result()
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "缓存错误"
		return resp, err
	}

	if exists > 0 {
		// 获取缓存中的 usercount
		userString, err := global.Redis.Get(ctx, key).Result()
		if err != nil {
			resp.StatusCode = 1
			resp.StatusMsg = "缓存错误"
			return resp, fmt.Errorf("usercount缓存错误：%v", err)
		}

		// 反序列化
		err = json.Unmarshal([]byte(userString), &usercount)
		if err != nil {
			resp.StatusCode = 1
			resp.StatusMsg = "缓存错误"
			return resp, err
		}
	} else {
		resp.StatusCode = 1
		resp.StatusMsg = "缓存错误"
		return resp, nil
	}

	// 更新 usercount 中的 total_favorited
	delta := int64(req.Delta) // 将 req.Delta 转换为 int64
	if req.Type == 1 {
		usercount.TotalFavorited += delta
	} else if req.Type == 2 {
		usercount.TotalFavorited -= delta
	}

	// 将更新后的 usercount 序列化并保存回 Redis
	updatedUserString, err := json.Marshal(usercount)
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "缓存错误"
		return resp, err
	}

	err = global.Redis.Set(ctx, key, updatedUserString, 0).Err()
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "缓存错误"
		return resp, err
	}

	// 操作成功
	resp.StatusCode = global.SuccessCode
	resp.StatusMsg = global.GetErrorMessage(global.SuccessCode)
	return resp, nil
}

func BuildUser(u *model.User, uc *model.UserCount) *user.User {
	userinfo := user.User{
		Id:              u.ID,
		Name:            u.Username,
		Avatar:          u.Avatar,
		BackgroundImage: u.BackgroundImage,
		Signature:       u.Signature,
		TotalFavorited:  uc.TotalFavorited,
		WorkCount:       uc.WorkCount,
		FavoriteCount:   uc.FavoriteCount,
		FollowCount:     uc.FollowCount,
		FollowerCount:   uc.FollowerCount,
	}
	return &userinfo
}
