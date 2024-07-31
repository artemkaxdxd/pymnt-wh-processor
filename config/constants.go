package config

import (
	"errors"
)

var ( // Errors
	ErrDuplicateEvent = errors.New("event already exists")
	ErrOrderEnded     = errors.New("order ended")
	ErrInvalidFilters = errors.New("invalid filters")
)

type ServiceCode int

const ( // WARNING: Do not insert constants in the middle, only add new to the end if needed
	CodeOK ServiceCode = iota
	CodeBadRequest
	CodeUnprocessableEntity
	CodeDatabaseError
	CodeNotFound
	CodeDuplicateEvent
	CodeOrderEnded
)

const ( // Messages
	MsgCreateOK = "create success"
	MsgUpdateOK = "update success"
)

type OrderStatus string

const (
	OrderStatusCreated     OrderStatus = "cool_order_created"
	OrderStatusPending     OrderStatus = "sbu_verification_pending"
	OrderStatusConfirmed   OrderStatus = "confirmed_by_mayor"
	OrderStatusChangedMind OrderStatus = "changed_my_mind"
	OrderStatusFailed      OrderStatus = "failed"
	OrderStatusChinazes    OrderStatus = "chinazes"
	OrderStatusMoneyBack   OrderStatus = "give_my_money_back"
)

var (
	OrderFinalStatuses = []OrderStatus{
		OrderStatusChangedMind, OrderStatusFailed,
		OrderStatusChinazes, OrderStatusMoneyBack,
	}

	OrderStatusesOrder = map[uint8]OrderStatus{
		1: OrderStatusCreated,
		2: OrderStatusPending,
		3: OrderStatusConfirmed,
		4: OrderStatusChangedMind,
		5: OrderStatusFailed,
		6: OrderStatusChinazes,
		7: OrderStatusMoneyBack,
	}
)
