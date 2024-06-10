package handler

import (
	"api/pb/social"
	"context"
	"social_service/global"
	"social_service/internal/model"
)

func (*SocialService) PostMessage(ctx context.Context, req *social.PostMessageRequest) (resp *social.PostMessageResponse, err error) {
	resp = new(social.PostMessageResponse)
	message := model.Message{
		UserId:   uint(req.UserId),
		ToUserId: uint(req.ToUserId),
		Content:  req.Content,
	}
	err = model.GetMessageInstance().CreateMessage(&message)
	if err != nil {
		resp.StatusCode = global.MYSQLQueryErrCode
		resp.StatusMsg = global.GetErrorMessage(global.MYSQLQueryErrCode)
		return resp, nil
	}
	resp.StatusCode = global.SuccessCode
	resp.StatusMsg = global.GetErrorMessage(global.SuccessCode)
	return resp, nil
}
func (*SocialService) GetMessage(ctx context.Context, req *social.GetMessageRequest) (resp *social.GetMessageResponse, err error) {
	resp = new(social.GetMessageResponse)
	var messages []*model.Message
	messages, _ = model.GetMessageInstance().GetMessageList(uint(req.UserId), uint(req.ToUserId), req.PreMsgTime)
	if err != nil {
		resp.StatusCode = global.MQSendErrCode
		resp.StatusMsg = global.GetErrorMessage(global.MQSendErrCode)
		return resp, nil

	}

	for _, message := range messages {
		resp.MessageList = append(resp.MessageList, &social.Message{
			Id:        int64(message.ID),
			UserId:    int64(message.UserId),
			ToUserId:  int64(message.ToUserId),
			Content:   message.Content,
			CreatedAt: message.CreatedAt.Unix(),
		})
	}
	resp.StatusCode = global.SuccessCode
	resp.StatusMsg = global.GetErrorMessage(global.SuccessCode)
	return resp, nil

}
