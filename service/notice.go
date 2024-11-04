package service

import (
	"context"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-telegram/bot"
	"github.com/siddontang/go-log/log"
	"telegram-robot/model"
)

type OrDerEventHandler struct {
	canal.DummyEventHandler
	TargetTable string
	Bot         *bot.Bot
	LastOrder   model.TgOrder // 存储上一次推送的数据
}

// OnRow 方法处理行事件，根据不同的操作类型进行处理
func (o *OrDerEventHandler) OnRow(e *canal.RowsEvent) error {
	// 获取发生变更的表名\
	fmt.Println(e.Table.Name, o.TargetTable)
	if e.Table.Name != o.TargetTable { // 只处理目标表的事件
		return nil
	}
	// 根据不同的操作类型进行响应处理
	switch e.Action {
	case canal.InsertAction:
		fmt.Println("ssss")
		o.handleInsert(e)
	case canal.UpdateAction:
		fmt.Println("xxx")
		o.handleUpdate(e)
	default:
		log.Warn("Unknown action detected")
	}
	return nil
}

func (o *OrDerEventHandler) handleInsert(e *canal.RowsEvent) {
	if len(e.Rows) > 0 {
		latestRow := e.Rows[len(e.Rows)-1] // 获取最新插入的行
		// 假设你已经将最新行映射到 TgOrder 结构体
		var t model.TgOrder
		if err := t.Find(); err != nil {
			fmt.Println(err)
			return
		}
		// 仅在数据有变更时才进行推送
		if o.isOrderChanged(t) {
			o.LastOrder = t                                  // 更新上一次插入的数据
			NotifyDataUpdate(context.Background(), o.Bot, t) // 调用推送函数
			log.Infof("Latest inserted row: %v", latestRow)
		}
	} else {
		log.Info("No rows found in insert event")
	}
}

// 处理更新操作的逻辑
func (h *OrDerEventHandler) handleUpdate(e *canal.RowsEvent) {
	if len(e.Rows) >= 2 {
		newRow := e.Rows[len(e.Rows)-1] // 新数据
		var t model.TgOrder
		if err := t.Find(); err != nil {
			fmt.Println(err)
			return
		}
		// 仅在数据有变更时才进行推送
		if h.isOrderChanged(t) {
			h.LastOrder = t                                  // 更新上一次推送的数据
			NotifyDataUpdate(context.Background(), h.Bot, t) // 调用推送函数
			log.Infof("Updated from %v to %v", e.Rows[len(e.Rows)-2], newRow)
		}
	}
}

// 检查订单是否发生变化
func (o *OrDerEventHandler) isOrderChanged(newOrder model.TgOrder) bool {
	return o.LastOrder.Phone != newOrder.Phone || o.LastOrder.QRCode != newOrder.QRCode || o.LastOrder.Code != newOrder.Code
}
