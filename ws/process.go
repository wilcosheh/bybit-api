package ws

func (b *ByBitWS) processOrderBookSnapshot(symbol string, ob ...*OrderBookL2) { // ob []*OrderBookL2
	var value *OrderBookLocal
	var ok bool

	value, ok = b.orderBookLocals[symbol]
	if !ok {
		value = NewOrderBookLocal()
		b.orderBookLocals[symbol] = value
	}
	value.LoadSnapshot(ob)

	b.Emit(WSOrderBook25L1, symbol, value.GetOrderBook())
}

func (b *ByBitWS) processOrderBookDelta(symbol string, delta *OrderBookL2Delta) {
	value, ok := b.orderBookLocals[symbol]
	if !ok {
		return
	}
	value.Update(delta)

	b.Emit(WSOrderBook25L1, symbol, value.GetOrderBook())
}

func (b *ByBitWS) processTrade(symbol string, data ...*Trade) {
	b.Emit(WSTrade, symbol, data)
}

func (b *ByBitWS) processKLine(symbol string, data KLine) {
	b.Emit(WSKLine, symbol, data)
}

func (b *ByBitWS) processKLineV2(symbol string, data []*KLineV2) {
	b.Emit(WSKLineV2, symbol, data)
}

func (b *ByBitWS) processCandle(symbol string, data []*KLineV2) {
	b.Emit(WSCandle, symbol, data)
}

func (b *ByBitWS) processInsurance(currency string, data ...*Insurance) {
	b.Emit(WSInsurance, currency, data)
}

func (b *ByBitWS) processInstrument(symbol string, data ...*Instrument) {
	b.Emit(WSInstrument, symbol, data)
}

func (b *ByBitWS) processLiquidation(symbol string, data *Liquidation) {
	b.Emit(WSLiquidation, symbol, data)
}

func (b *ByBitWS) processPosition(data ...*Position) {
	b.Emit(WSPosition, data)
}

func (b *ByBitWS) processExecution(data ...*Execution) {
	b.Emit(WSExecution, data)
}

func (b *ByBitWS) processOrder(data ...*Order) {
	b.Emit(WSOrder, data)
}

func (b *ByBitWS) processStopOrder(data ...*StopOrder) {
	b.Emit(WSStopOrder, data)
}

func (b *ByBitWS) processWallet(data ...*Wallet) {
	b.Emit(WSWallet, data)
}
