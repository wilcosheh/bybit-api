package ws

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/chuckpreslar/emission"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"github.com/wilcosheh/bybit-api/recws"
)

const (
	MaxTryTimes = 10
)

// https://github.com/bybit-exchange/bybit-official-api-docs/blob/master/zh_cn/websocket.md

// 测试网地址
// wss://stream-testnet.bybit.com/realtime_public
// wss://stream-testnet.bybit.com/realtime_private

// 主网地址
// wss://stream.bybit.com/realtime_public
// wss://stream.bybit.com/realtime_private

const (
	HostRealPublic     = "wss://stream.bybit.com/realtime_public"
	HostRealPrivate    = "wss://stream.bybit.com/realtime_private"
	HostTestnetPublic  = "wss://stream-testnet.bybit.com/realtime_public"
	HostTestnetPrivate = "wss://stream-testnet.bybit.com/realtime_private"
)

const (
	WSOrderBook25L1 = "orderBookL2_25" // 新版25档orderBook: order_book_25L1.BTCUSD
	WSKLine         = "kline"          // K线: kline.BTCUSD.1m
	WSKLineV2       = "klineV2"        // V2版本K线: klineV2.1.BTCUSD
	WSCandle        = "candle"         // USDT永续K线: candle.1.BTCUSDT
	WSTrade         = "trade"          // 实时交易: trade/trade.BTCUSD
	WSInsurance     = "insurance"      // 每日保险基金更新: insurance
	WSInstrument    = "instrument"     // 产品最新行情: instrument
	WSLiquidation   = "liquidation"    // 強平推送: liquidation

	WSPosition  = "position"   // 仓位变化: position
	WSExecution = "execution"  // 委托单成交信息: execution
	WSOrder     = "order"      // 委托单的更新: order
	WSStopOrder = "stop_order" // 条件单的更新: stop_order
	WSWallet    = "wallet"     // 条件单的更新: stop_order

	WSDisconnected = "disconnected" // WS断开事件
)

var (
	topicOrderBook25l1prefix = WSOrderBook25L1 + "."
)

type Configuration struct {
	Addr          string `json:"addr"`
	Proxy         string `json:"proxy"` // http://127.0.0.1:1081
	ApiKey        string `json:"api_key"`
	SecretKey     string `json:"secret_key"`
	AutoReconnect bool   `json:"auto_reconnect"`
	DebugMode     bool   `json:"debug_mode"`
}

type ByBitWS struct {
	cfg    *Configuration
	ctx    context.Context
	cancel context.CancelFunc
	conn   *recws.RecConn
	mu     sync.RWMutex
	Ended  bool

	subscribeCmds   []Cmd
	orderBookLocals map[string]*OrderBookLocal // key: symbol

	emitter *emission.Emitter
}

func New(config *Configuration) *ByBitWS {
	b := &ByBitWS{
		cfg:             config,
		emitter:         emission.NewEmitter(),
		orderBookLocals: make(map[string]*OrderBookLocal),
	}
	b.ctx, b.cancel = context.WithCancel(context.Background())

	b.conn = &recws.RecConn{
		KeepAliveTimeout: 60 * time.Second,
		NonVerbose:       true,
	}
	if config.Proxy != "" {
		proxy, err := url.Parse(config.Proxy)
		if err != nil {
			return nil
		}
		b.conn.Proxy = http.ProxyURL(proxy)
	}
	b.conn.SubscribeHandler = b.subscribeHandler
	return b
}

func (b *ByBitWS) subscribeHandler() error {
	if b.cfg.DebugMode {
		log.Printf("BybitWs subscribeHandler")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.cfg.ApiKey != "" && b.cfg.SecretKey != "" {
		err := b.Auth()
		if err != nil {
			log.Printf("BybitWs auth error: %v", err)
		}
	}

	for _, cmd := range b.subscribeCmds {
		err := b.SendCmd(cmd)
		if err != nil {
			log.Printf("BybitWs SendCmd return error: %v", err)
		}
	}

	return nil
}

func (b *ByBitWS) closeHandler(code int, text string) error {
	if b.cfg.DebugMode {
		log.Printf("BybitWs close handle executed code=%v text=%v", code, text)
	}
	return nil
}

// IsConnected returns the WebSocket connection state
func (b *ByBitWS) IsConnected() bool {
	return b.conn.IsConnected()
}

func (b *ByBitWS) Subscribe(arg string) {
	cmd := Cmd{
		Op:   "subscribe",
		Args: []interface{}{arg},
	}
	b.subscribeCmds = append(b.subscribeCmds, cmd)
	b.SendCmd(cmd)
}

func (b *ByBitWS) SendCmd(cmd Cmd) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return b.Send(string(data))
}

func (b *ByBitWS) Send(msg string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("BybitWs send error: %v", r)
		}
	}()

	err = b.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	return
}

func (b *ByBitWS) Start() error {
	b.connect()

	cancel := make(chan struct{})

	go func() {
		t := time.NewTicker(time.Second * 5)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				b.ping()
			case <-cancel:
				return
			}
		}
	}()

	go func() {
		defer close(cancel)

		for {
			messageType, data, err := b.conn.ReadMessage()
			if err != nil {
				log.Printf("BybitWs Read error, closing connection: %v", err)
				b.conn.Close()
				b.Ended = true
				return
			}

			if err = b.processMessage(messageType, data); err != nil {
				log.Printf("BybitWs process error: %v", err)
			}
		}
	}()

	return nil
}

func (b *ByBitWS) connect() {
	b.conn.Dial(b.cfg.Addr, nil)
}

func (b *ByBitWS) ping() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("BybitWs ping error: %v", r)
		}
	}()

	if !b.IsConnected() {
		return
	}
	err := b.conn.WriteMessage(websocket.TextMessage, []byte(`{"op":"ping"}`))
	if err != nil {
		log.Printf("BybitWs ping error: %v", err)
	}
}

func (b *ByBitWS) Auth() error {
	// 单位:毫秒
	expires := time.Now().Unix()*1000 + 10000
	req := fmt.Sprintf("GET/realtime%d", expires)
	sig := hmac.New(sha256.New, []byte(b.cfg.SecretKey))
	sig.Write([]byte(req))
	signature := hex.EncodeToString(sig.Sum(nil))

	cmd := Cmd{
		Op: "auth",
		Args: []interface{}{
			b.cfg.ApiKey,
			//fmt.Sprintf("%v", expires),
			expires,
			signature,
		},
	}
	err := b.SendCmd(cmd)
	return err
}

func (b *ByBitWS) processMessage(messageType int, data []byte) error {
	ret := gjson.ParseBytes(data)

	if b.cfg.DebugMode {
		log.Printf("BybitWs %v", string(data))
	}

	// 处理心跳包
	if ret.Get("ret_msg").String() == "pong" {
		b.handlePong()
	}

	if topicValue := ret.Get("topic"); topicValue.Exists() {
		topic := topicValue.String()
		if strings.HasPrefix(topic, topicOrderBook25l1prefix) {
			symbol := topic[len(topicOrderBook25l1prefix):]
			type_ := ret.Get("type").String()
			raw := ret.Get("data").Raw

			switch type_ {
			case "snapshot":
				var data []*OrderBookL2
				raw = ret.Get("data.order_book").Raw
				err := json.Unmarshal([]byte(raw), &data)
				if err != nil {
					return err
				}
				b.processOrderBookSnapshot(symbol, data...)
			case "delta":
				var delta OrderBookL2Delta
				err := json.Unmarshal([]byte(raw), &delta)
				if err != nil {
					return err
				}
				b.processOrderBookDelta(symbol, &delta)
			}
		} else if strings.HasPrefix(topic, WSTrade) {
			symbol := strings.TrimLeft(topic, WSTrade+".")
			raw := ret.Get("data").Raw
			var data []*Trade
			err := json.Unmarshal([]byte(raw), &data)
			if err != nil {
				return err
			}
			b.processTrade(symbol, data...)
		} else if strings.HasPrefix(topic, WSCandle) {
			// candle.1.BTCUSDT
			topicArray := strings.Split(topic, ".")
			if len(topicArray) != 3 {
				return errors.New("candle topic format error")
			}
			symbol := topicArray[2]
			raw := ret.Get("data").Raw
			var data []*KLineV2
			err := json.Unmarshal([]byte(raw), &data)
			if err != nil {
				return err
			}
			for _, kline := range data {
				kline.Symbol = symbol
				kline.Interval = topicArray[1]
			}
			b.processCandle(symbol, data)
		} else if strings.HasPrefix(topic, WSKLineV2) {
			// klineV2.1.BTCUSD
			topicArray := strings.Split(topic, ".")
			if len(topicArray) != 3 {
				return errors.New("klineV2 topic format error")
			}
			symbol := topicArray[2]
			raw := ret.Get("data").Raw
			var data []*KLineV2
			err := json.Unmarshal([]byte(raw), &data)
			if err != nil {
				return err
			}
			for _, kline := range data {
				kline.Symbol = symbol
				kline.Interval = topicArray[1]
			}
			b.processKLineV2(symbol, data)
		} else if strings.HasPrefix(topic, WSKLine) {
			// kline.BTCUSD.1m
			topicArray := strings.Split(topic, ".")
			if len(topicArray) != 3 {
				return errors.New("kline topic format error")
			}
			symbol := topicArray[1]
			raw := ret.Get("data").Raw
			var data KLine
			err := json.Unmarshal([]byte(raw), &data)
			if err != nil {
				return err
			}
			b.processKLine(symbol, data)
		} else if strings.HasPrefix(topic, WSInsurance) {
			// insurance.BTC
			topicArray := strings.Split(topic, ".")
			if len(topicArray) != 2 {
				return errors.New("insurance topic format error")
			}
			currency := topicArray[1]
			raw := ret.Get("data").Raw
			var data []*Insurance
			err := json.Unmarshal([]byte(raw), &data)
			if err != nil {
				return err
			}
			b.processInsurance(currency, data...)
		} else if strings.HasPrefix(topic, WSInstrument) {
			topicArray := strings.Split(topic, ".")
			if len(topicArray) != 2 {
				return errors.New("instrument topic format error")
			}
			symbol := topicArray[1]
			raw := ret.Get("data").Raw
			var data []*Instrument
			err := json.Unmarshal([]byte(raw), &data)
			if err != nil {
				return err
			}
			b.processInstrument(symbol, data...)
		} else if strings.HasPrefix(topic, WSLiquidation) {
			topicArray := strings.Split(topic, ".")
			if len(topicArray) != 2 {
				return errors.New("liquidation topic format error")
			}
			symbol := topicArray[1]
			raw := ret.Get("data").Raw
			var data *Liquidation
			err := json.Unmarshal([]byte(raw), &data)
			if err != nil {
				return err
			}
			b.processLiquidation(symbol, data)
		} else if topic == WSPosition {
			raw := ret.Get("data").Raw
			var data []*Position
			err := json.Unmarshal([]byte(raw), &data)
			if err != nil {
				return err
			}
			b.processPosition(data...)
		} else if topic == WSExecution {
			raw := ret.Get("data").Raw
			var data []*Execution
			err := json.Unmarshal([]byte(raw), &data)
			if err != nil {
				return err
			}
			b.processExecution(data...)
		} else if topic == WSOrder {
			raw := ret.Get("data").Raw
			var data []*Order
			err := json.Unmarshal([]byte(raw), &data)
			if err != nil {
				return err
			}

			for _, order := range data {
				if order.Timestamp.IsZero() && !order.CreateTime.IsZero() {
					order.Timestamp = order.CreateTime
				}
				if !order.Timestamp.IsZero() && order.CreateTime.IsZero() {
					order.CreateTime = order.Timestamp
				}
			}
			b.processOrder(data...)
		} else if topic == WSStopOrder {
			raw := ret.Get("data").Raw
			var data []*StopOrder
			err := json.Unmarshal([]byte(raw), &data)
			if err != nil {
				return err
			}
			b.processStopOrder(data...)
		} else if topic == WSWallet {
			raw := ret.Get("data").Raw
			var data []*Wallet
			err := json.Unmarshal([]byte(raw), &data)
			if err != nil {
				return err
			}
			b.processWallet(data...)
		}
		return nil
	}

	return nil
}

func (b *ByBitWS) handlePong() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("handlePong error: %v", r)
		}
	}()
	pongHandler := b.conn.PongHandler()
	if pongHandler != nil {
		pongHandler("pong")
	}
	return nil
}

func (b *ByBitWS) CloseAndReconnect() {
	b.conn.CloseAndReconnect()
}
