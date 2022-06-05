package ws

import (
	sjson "encoding/json"
	"strconv"
	"time"
)

type OrderBookL2 struct {
	ID     int64   `json:"id,string"`
	Price  float64 `json:"price,string"`
	Side   string  `json:"side"`
	Size   float64 `json:"size"`
	Symbol string  `json:"symbol"`
}

type OrderBookL2Delta struct {
	Delete []*OrderBookL2 `json:"delete"`
	Update []*OrderBookL2 `json:"update"`
	Insert []*OrderBookL2 `json:"insert"`
}

func (o *OrderBookL2) Key() string {
	return strconv.FormatInt(o.ID, 10)
}

type Trade struct {
	Timestamp     time.Time `json:"timestamp"`
	Symbol        string    `json:"symbol"`
	Side          string    `json:"side"`
	Size          float64   `json:"size"`
	Price         float64   `json:"price"`
	TickDirection string    `json:"tick_direction"`
	TradeID       string    `json:"trade_id"`
	CrossSeq      int       `json:"cross_seq"` // only valid for inverse
}

type KLine struct {
	ID       int64   `json:"id"`        // 563
	Symbol   string  `json:"symbol"`    // BTCUSD
	OpenTime int64   `json:"open_time"` // 1539918000
	Open     float64 `json:"open"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Close    float64 `json:"close"`
	Volume   float64 `json:"volume"`
	Turnover float64 `json:"turnover"` // 0.0013844
	Interval string  `json:"interval"` // 1m
}

type KLineV2 struct {
	Symbol    string  `json:"symbol"`          // 合约类型，从 topic 解析得到
	Start     int64   `json:"start"`           // 开始时间戳（秒）
	End       int64   `json:"end"`             // 结束时间戳（秒）
	Open      float64 `json:"open"`            // 开盘价
	Close     float64 `json:"close"`           // 收盘价
	High      float64 `json:"high"`            // 最高价格
	Low       float64 `json:"low"`             // 最低价格
	Volume    float64 `json:"volume,string"`   // 交易量 TODO: 反向永续类型不一样
	Turnover  float64 `json:"turnover,string"` // 成交金额 0.0013844 TODO: 反向永续类型不一样
	Confirm   bool    `json:"confirm"`         // 是否确认，为 true 表明是 k 线 最后一个 tick，否则只是一个快照数据，即中间价格
	CrossSeq  int     `json:"cross_seq"`       // 版本号
	Interval  string  `json:"interval"`        // 周期，从 topic 解析得到： 1 3 5 15 30 60 120 240 360 D W M
	Timestamp int64   `json:"timestamp"`       // 结束时间戳（秒）
}

type Insurance struct {
	Currency      string    `json:"currency"`
	Timestamp     time.Time `json:"timestamp"`
	WalletBalance int64     `json:"wallet_balance"`
}

type Instrument struct {
	Symbol     string  `json:"symbol"`
	MarkPrice  float64 `json:"mark_price"`
	IndexPrice float64 `json:"index_price"`
}

type Liquidation struct {
	Symbol string       `json:"symbol"` // 合约类型
	Side   string       `json:"side"`   // 被强平仓位的方向
	Price  sjson.Number `json:"price"`  // 破产价格
	Qty    float64      `json:"qty"`    // 交易數量
	Time   int64        `json:"time"`   // 毫秒時間戳
}

type Order struct {
	OrderID        string       `json:"order_id"`            // 订单ID
	OrderLinkID    string       `json:"order_link_id"`       // 自定义订单ID
	Symbol         string       `json:"symbol"`              // 合约类型
	Side           string       `json:"side"`                // 方向
	OrderType      string       `json:"order_type"`          // 委托单价格类型，Limit/Market
	Price          sjson.Number `json:"price"`               // 委托价格
	Qty            float64      `json:"qty"`                 // 委托数量
	TimeInForce    string       `json:"time_in_force"`       // 执行策略，GoodTillCancel/ImmediateOrCancel/FillOrKill/PostOnly
	CreateType     string       `json:"create_type"`         // 下单操作的触发场景
	CancelType     string       `json:"cancel_type"`         // 取消操作的触发场景
	OrderStatus    string       `json:"order_status"`        // 订单状态
	LeavesQty      float64      `json:"leaves_qty"`          // 剩余委托数量
	CumExecQty     float64      `json:"cum_exec_qty"`        // 累计成交数量
	CumExecValue   sjson.Number `json:"cum_exec_value"`      // 累计成交价值
	CumExecFee     sjson.Number `json:"cum_exec_fee"`        // 累计成交手续费
	Timestamp      time.Time    `json:"timestamp"`           // 创建时间，only valid for inverse
	CreateTime     time.Time    `json:"create_time"`         // 创建时间，only valid for linear
	UpdateTime     time.Time    `json:"update_time"`         // 成交时间，only valid for linear
	TakeProfit     sjson.Number `json:"take_profit"`         // 止盈价格
	StopLoss       sjson.Number `json:"stop_loss"`           // 止损价格
	TrailingStop   sjson.Number `json:"trailing_stop"`       // 追踪止损（与当前价格的距离）
	TrailingActive sjson.Number `json:"trailing_active"`     // 激活价格
	LastExecPrice  sjson.Number `json:"last_exec_price"`     // 最近一次成交价格
	ReduceOnly     bool         `json:"reduce_only"`         // 只减仓
	PositionIdx    int          `json:"position_idx,string"` // 用于在不同仓位模式下标识仓位：0 - 单向持仓，1 - 双向持仓Buy，2 - 双向持仓Sell，only valid for linear
	CloseOnTrigger bool         `json:"close_on_trigger"`    // 触发后平仓，如果下平仓单，请设置为 true，避免因为保证金不足而导致下单失败
}

type StopOrder struct {
	OrderID        string       `json:"order_id"`
	OrderLinkID    string       `json:"order_link_id"`
	UserID         int64        `json:"user_id"`
	Symbol         string       `json:"symbol"`
	Side           string       `json:"side"`
	OrderType      string       `json:"order_type"`
	Price          sjson.Number `json:"price"`
	Qty            float64      `json:"qty"`
	TimeInForce    string       `json:"time_in_force"` // GoodTillCancel/ImmediateOrCancel/FillOrKill/PostOnly
	CreateType     string       `json:"create_type"`
	CancelType     string       `json:"cancel_type"`
	OrderStatus    string       `json:"order_status"`
	StopOrderType  string       `json:"stop_order_type"`
	TriggerBy      string       `json:"trigger_by"`
	TriggerPrice   sjson.Number `json:"trigger_pricee"`
	CloseOnTrigger bool         `json:"close_on_trigger"`
	Timestamp      time.Time    `json:"timestamp"`
}

type Execution struct {
	Symbol      string    `json:"symbol"`        // 合约类型
	Side        string    `json:"side"`          // 方向
	OrderID     string    `json:"order_id"`      // 订单ID
	ExecID      string    `json:"exec_id"`       // 成交ID
	OrderLinkID string    `json:"order_link_id"` // 自定义订单ID
	Price       float64   `json:"price"`         // 成交价格
	OrderQty    float64   `json:"order_qty"`     // 订单数量
	ExecType    string    `json:"exec_type"`     // 交易类型，Trade/AdlTrade/BustTrade
	ExecQty     float64   `json:"exec_qty"`      // 成交数量
	ExecFee     float64   `json:"exec_fee"`      // 交易手续费
	LeavesQty   float64   `json:"leaves_qty"`    // 剩余委托数量
	IsMaker     bool      `json:"is_maker"`      // 是否是maker
	TradeTime   time.Time `json:"trade_time"`    // 交易时间
}

type Position struct {
	UserID         int64   `json:"user_id,string"`      // 用户 ID
	Symbol         string  `json:"symbol"`              // 合约类型
	Size           float64 `json:"size"`                // 仓位数量
	Side           string  `json:"side"`                // 方向
	PositionValue  float64 `json:"position_value"`      // 仓位价值
	EntryPrice     float64 `json:"entry_price"`         // 平均入场价
	LiqPrice       float64 `json:"liq_price"`           // 强平价格
	BustPrice      float64 `json:"bust_price"`          // 破产价格
	Leverage       float64 `json:"leverage"`            // 逐仓模式下，值为一哦哪个好设置的杠杆；全仓模式下，值为当前风险限额下最大杠杆
	OrderMargin    float64 `json:"order_margin"`        // 委托预占用保证金
	PositionMargin float64 `json:"position_margin"`     // 仓位保证金
	OccClosingFee  float64 `json:"occ_closing_fee"`     // 仓位占用的平仓手续费
	TakeProfit     float64 `json:"take_profit"`         // 止盈价格
	TpTriggerBy    string  `json:"tp_trigger_by"`       // 止盈激活价格类型，默认为 LastPrice
	StopLoss       float64 `json:"stop_loss"`           // 止损价格
	SlTriggerBy    string  `json:"sl_trigger_by"`       // 止损激活价格类型
	RealisedPnl    float64 `json:"realised_pnl"`        // 当日已结盈亏
	CumRealisedPnl float64 `json:"cum_realised_pnl"`    // 累计已结盈亏
	PositionStatus string  `json:"position_status"`     // 仓位状态：正常、强平、减仓
	PositionSeq    int64   `json:"position_seq,string"` // 仓位变化版本号
	PositionIdx    int     `json:"position_idx,string"` // 用于在不同仓位模式下标识仓位：0 - 单向持仓，1 - 双向持仓Buy，2 - 双向持仓Sell，only valid for linear
	Mode           string  `json:"mode"`                // 仓位模式： MergedSingle or BothSide
	Isolated       bool    `json:"isolated"`            // 是否逐仓，true-逐仓 false-全仓
	RiskID         int     `json:"risk_id,string"`      // 风险限额 ID

	// 反向永续字段
	TrailingStop     float64 `json:"trailing_stop"` //
	TrailingActive   float64 `json:"trailing_active"`
	WalletBalance    float64 `json:"wallet_balance"`
	AvailableBalance float64 `json:"available_balance"`
	OccFundingFee    float64 `json:"occ_funding_fee"`
	AutoAddMargin    int     `json:"auto_add_margin,string"`
}

type Wallet struct {
	WalletBalance    float64 `json:"wallet_balance"`
	AvailableBalance float64 `json:"available_balance"`
}
