package model

import "time"

/*
用户订单结构实体定义
*/
type UserOrder struct {
	Id            int64        `xorm:"pk autoincr" json:"id"` //主键 自增
	SumMoney      int64        `xorm:"default 0" json:"sum_money"`
	Time          time.Time    `xorm:"DateTime" json:"time"`
	OrderTime     uint64       `json:"order_time"`
	OrderStatusId int64        `xorm:"index" json:"order_status_id"`
	OrderStatus   *OrderStatus `xorm:"-"`
	UserId        int64        `xorm:"index" json:"user_id"`
	User          *User        `xorm:"-"`
	ShopId        int64        `xorm:"index" json:"shop_id"`
	Shop          *Shop        `xorm:"-"`
	AddressId     int64        `xorm:"index" json:"address_id"`
	Address       *Address     `xorm:"-"`
	DelFlag       int64        `xorm:"default 0" json:"del_flag"`
}
