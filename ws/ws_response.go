package ws

import (
	sjson "encoding/json"
	"strconv"
	"time"
)

type OrderBookL2 struct {
	ID     int64   `json:"id"`
	Price  float64 `json:"price,string"`
	Side   string  `json:"side"`
	Size   int64   `json:"size"`
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
	Size          int       `json:"size"`
	Price         float64   `json:"price"`
	TickDirection string    `json:"tick_direction"`
	TradeID       string    `json:"trade_id"`
	CrossSeq      int       `json:"cross_seq"`
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
	Symbol    string  `json:"symbol"`   // BTCUSD
	Start     int64   `json:"start"`    // start time of the candle 1572425640
	End       int64   `json:"end"`      // end time of the candle 1572425700
	Open      float64 `json:"open"`     // open price
	Close     float64 `json:"close"`    // close price
	High      float64 `json:"high"`     // max price
	Low       float64 `json:"low"`      // min price
	Volume    float64 `json:"volume"`   // volume 81790
	Turnover  float64 `json:"turnover"` // turnover 8.889247899999999
	Confirm   bool    `json:"confirm"`  // snapshot flag
	CrossSeq  int     `json:"cross_seq"`
	Timestamp int64   `json:"timestamp"` // cross time 1572425676958323
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
	OrderID        string       `json:"order_id"`
	OrderLinkID    string       `json:"order_link_id"`
	Symbol         string       `json:"symbol"`
	Side           string       `json:"side"`
	OrderType      string       `json:"order_type"`
	Price          sjson.Number `json:"price"`
	Qty            float64      `json:"qty"`
	TimeInForce    string       `json:"time_in_force"` // GoodTillCancel/ImmediateOrCancel/FillOrKill/PostOnly
	CreateType     string       `json:"create_type"`
	CancelType     string       `json:"cancel_type"`
	OrderStatus    string       `json:"order_status"`
	LeavesQty      float64      `json:"leaves_qty"`
	CumExecQty     float64      `json:"cum_exec_qty"`
	CumExecValue   sjson.Number `json:"cum_exec_value"`
	CumExecFee     sjson.Number `json:"cum_exec_fee"`
	Timestamp      time.Time    `json:"timestamp"`
	TakeProfit     sjson.Number `json:"take_profit"`
	StopLoss       sjson.Number `json:"stop_loss"`
	TrailingStop   sjson.Number `json:"trailing_stop"`
	TrailingActive sjson.Number `json:"trailing_active"`
	LastExecPrice  sjson.Number `json:"last_exec_price"`
	ReduceOnly     bool         `json:"reduce_only"`
	CloseOnTrigger bool         `json:"close_on_trigger"`
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
	Symbol      string    `json:"symbol"`
	Side        string    `json:"side"`
	OrderID     string    `json:"order_id"`
	ExecID      string    `json:"exec_id"`
	OrderLinkID string    `json:"order_link_id"`
	Price       float64   `json:"price,string"`
	OrderQty    float64   `json:"order_qty"`
	ExecType    string    `json:"exec_type"`
	ExecQty     float64   `json:"exec_qty"`
	ExecFee     float64   `json:"exec_fee,string"`
	LeavesQty   float64   `json:"leaves_qty"`
	IsMaker     bool      `json:"is_maker"`
	TradeTime   time.Time `json:"trade_time"`
}

type Position struct {
	UserID           int64   `json:"user_id"`
	Symbol           string  `json:"symbol"`
	Size             float64 `json:"size"`
	Side             string  `json:"side"`
	PositionValue    float64 `json:"position_value,string"`
	EntryPrice       float64 `json:"entry_price,string"`
	LiqPrice         float64 `json:"liq_price,string"`
	BustPrice        float64 `json:"bust_price,string"`
	Leverage         float64 `json:"leverage,string"`
	OrderMargin      float64 `json:"order_margin,string"`
	PositionMargin   float64 `json:"position_margin,string"`
	AvailableBalance float64 `json:"available_balance,string"`
	TakeProfit       float64 `json:"take_profit,string"`
	TpTriggerBy      string  `json:"tp_trigger_by"`
	StopLoss         float64 `json:"stop_loss,string"`
	SlTriggerBy      string  `json:"sl_trigger_by"`
	RealisedPnl      float64 `json:"realised_pnl,string"`
	TrailingStop     float64 `json:"trailing_stop,string"`
	TrailingActive   float64 `json:"trailing_active,string"`
	WalletBalance    float64 `json:"wallet_balance,string"`
	RiskID           int     `json:"risk_id"`
	OccClosingFee    float64 `json:"occ_closing_fee,string"`
	OccFundingFee    float64 `json:"occ_funding_fee,string"`
	AutoAddMargin    int     `json:"auto_add_margin"`
	CumRealisedPnl   float64 `json:"cum_realised_pnl,string"`
	PositionStatus   string  `json:"position_status"`
	PositionSeq      int64   `json:"position_seq"`
}
