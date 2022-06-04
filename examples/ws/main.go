package main

import (
	"log"

	"github.com/wilcosheh/bybit-api/ws"
)

func main() {

	wsPublic := ws.New(&ws.Configuration{
		Addr:          ws.HostRealPublic, // 测试网络
		AutoReconnect: true,              // 断线自动重连
		DebugMode:     true,
	})
	wsPrivate := ws.New(&ws.Configuration{
		Addr:          ws.HostRealPrivate, // 测试网络
		ApiKey:        "wKuYtkeNdC2PaMKjoy",
		SecretKey:     "5ekcDn3KnKoCRbfvrPImYzVdx7Ri2hhVxkmw",
		AutoReconnect: true, // 断线自动重连
		DebugMode:     true,
	})

	// 订阅新版25档orderBook
	wsPublic.Subscribe(ws.WSOrderBook25L1 + ".ETHUSDT")
	// 实时交易
	wsPublic.Subscribe("trade.BTCUSD")
	wsPublic.Subscribe(ws.WSTrade) // BTCUSD/ETHUSD/EOSUSD/XRPUSD
	// K线
	wsPublic.Subscribe(ws.WSKLineV2 + ".1.BTCUSD")
	// 每日保险基金更新
	wsPublic.Subscribe(ws.WSInsurance)
	// 产品最新行情
	wsPublic.Subscribe(ws.WSInstrument + ".BTCUSD")

	// 仓位变化
	wsPrivate.Subscribe(ws.WSPosition)
	// 委托单成交信息
	wsPrivate.Subscribe(ws.WSExecution)
	// 委托单的更新
	wsPrivate.Subscribe(ws.WSOrder)

	wsPublic.On(ws.WSOrderBook25L1, handleOrderBook)
	wsPublic.On(ws.WSTrade, handleTrade)
	wsPublic.On(ws.WSKLineV2, handleKLineV2)
	wsPublic.On(ws.WSInsurance, handleInsurance)
	wsPublic.On(ws.WSInstrument, handleInstrument)

	wsPrivate.On(ws.WSPosition, handlePosition)
	wsPrivate.On(ws.WSExecution, handleExecution)
	wsPrivate.On(ws.WSOrder, handleOrder)

	wsPublic.Start()
	wsPrivate.Start()

	forever := make(chan struct{})
	<-forever
}

func handleOrderBook(symbol string, data ws.OrderBook) {
	log.Printf("handleOrderBook %v/%v", symbol, data)
}

func handleTrade(symbol string, data []*ws.Trade) {
	log.Printf("handleTrade %v/%v", symbol, data)
}

func handleKLineV2(symbol string, data []*ws.KLineV2) {
	for _, kline := range data {
		log.Printf("handleKLine %v/%v", symbol, kline)
	}
}

func handleInsurance(currency string, data []*ws.Insurance) {
	log.Printf("handleInsurance %v/%v", currency, data)
}

func handleInstrument(symbol string, data []*ws.Instrument) {
	log.Printf("handleInstrument %v/%v", symbol, data)
}

func handlePosition(data []*ws.Position) {
	log.Printf("handlePosition %v", data)
}

func handleExecution(data []*ws.Execution) {
	log.Printf("handleExecution %v", data)
}

func handleOrder(data []*ws.Order) {
	log.Printf("handleOrder %v", data)
}
