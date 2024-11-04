package main

import (
	"context"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/siddontang/go-log/log"
	"os"
	"os/signal"
	"telegram-robot/initialize/config"
	"telegram-robot/initialize/mysql"
	"telegram-robot/service"
)

func main() {
	config.InitConfig() //初始配置文件
	mysql.InitMysql()   // 初始化 MySQL

	// 初始化 Canal 配置
	cfg := canal.NewDefaultConfig()
	cfg.Addr = config.Conf.Mysql.Host + ":" + config.Conf.Mysql.Port
	cfg.User = config.Conf.Mysql.UserName
	cfg.Password = config.Conf.Mysql.Password // 如果有密码，设置密码
	cfg.Dump.TableDB = config.Conf.Mysql.DriverName
	cfg.Dump.Tables = []string{""} // 监听的表
	cfg.Dump.ExecutionPath = ""    //设置为空
	fmt.Println(config.Conf.Mysql.TargetTable)
	// 创建 Canal 实例
	c, err := canal.NewCanal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// 启动 Telegram Bot
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	opts := []bot.Option{
		bot.WithDefaultHandler(handler), //监听函数(自动回复)
	}
	b, err := bot.New(config.Conf.Token, opts...)
	if err != nil {
		panic(err)
	}
	fmt.Println("机器人已经启动")

	targetTable := ""                                                               // 注册自定义事件处理器，并设置目标表
	c.SetEventHandler(&service.OrDerEventHandler{TargetTable: targetTable, Bot: b}) //监听函数(自定义操作)

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
