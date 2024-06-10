package handler

import (
	"api/pb/social"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"social_service/internal/model"
	"video_service/global"
)

type SocialService struct {
	social.UnimplementedSocialServiceServer // 版本兼容问题
}

func NewSocialService() *SocialService {
	return &SocialService{}
}

// FollowAction 关注操作
func (*SocialService) FollowAction(ctx context.Context, req *social.FollowRequest) (resp *social.FollowResponse, err error) {
	actionType := req.ActionType
	relationModel := model.GetRelationInstance()
	if actionType == 1 {
		// 关注操作
		err = relationModel.CreateRelation(&model.Relation{
			UserId:   uint(req.UserId),
			TargetId: uint(req.ToUserId),
		})
		if err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "Failed to follow user"
			return resp, err
		}
		resp.StatusCode = 0
		resp.StatusMsg = "Successfully followed user"
	} else if actionType == 2 {
		// 取消关注操作
		err = relationModel.DeleteRelation(uint(req.UserId), uint(req.ToUserId))
		if err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "Failed to unfollow user"
			return resp, err
		}
		resp.StatusCode = 0
		resp.StatusMsg = "Successfully unfollowed user"
	} else {
		resp.StatusCode = -1
		resp.StatusMsg = "Invalid action type"
		return resp, status.Errorf(codes.InvalidArgument, "Invalid action type")
	}

	return resp, nil
}

func (*SocialService) GetFollowList(ctx context.Context, req *social.FollowListRequest) (resp *social.FollowListResponse, err error) {
	resp = new(social.FollowListResponse)
	resp.UserId, err = model.GetRelationInstance().GetFollowList(uint(req.UserId))
	if err != nil {
		resp.StatusCode = global.MYSQLQueryErrCode
		resp.StatusMsg = global.GetErrorMessage(global.MYSQLQueryErrCode)
		return resp, err
	}
	resp.StatusCode = 0
	resp.StatusMsg = "Successfully get follow list"
	return resp, nil
}

func (*SocialService) GetFollowerList(ctx context.Context, req *social.FollowListRequest) (resp *social.FollowListResponse, err error) {
	resp = new(social.FollowListResponse)
	resp.UserId, err = model.GetRelationInstance().GetFollowerList(uint(req.UserId))
	if err != nil {
		resp.StatusCode = global.MYSQLQueryErrCode
		resp.StatusMsg = global.GetErrorMessage(global.MYSQLQueryErrCode)
		return resp, err
	}
	resp.StatusCode = 0
	resp.StatusMsg = "Successfully get follow list"
	return resp, nil
}

func (*SocialService) GetFriendList(ctx context.Context, req *social.FollowListRequest) (resp *social.FollowListResponse, err error) {
	resp = new(social.FollowListResponse)
	resp.UserId, err = model.GetRelationInstance().GetFriendList(uint(req.UserId))
	if err != nil {
		resp.StatusCode = global.MYSQLQueryErrCode
		resp.StatusMsg = global.GetErrorMessage(global.MYSQLQueryErrCode)
		return resp, err
	}
	resp.StatusCode = 0
	resp.StatusMsg = "Successfully get follow list"
	return resp, nil
}

func (*SocialService) GetFollowInfo(ctx context.Context, req *social.FollowInfoRequest) (resp *social.FollowInfoResponse, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFollowInfo not implemented")
}
