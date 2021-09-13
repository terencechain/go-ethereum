package txpool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

type dumbTxPool struct {
	gasPrice *big.Int

	pendingNonces *txNoncer

	txFeed       event.Feed
	scope        event.SubscriptionScope
	chainHeadSub event.Subscription
}

func (pool *dumbTxPool) SetGasPrice(price *big.Int) {
	pool.gasPrice = price
	// TODO reorg all tx, remove txs below gp
}

func (pool *dumbTxPool) Nonce(addr common.Address) uint64 {
	return pool.pendingNonces.get(addr)
}

func (pool *dumbTxPool) Stats() (int, int) {
	return 0, 0
}

func (pool *dumbTxPool) Content() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	// Only needed in miner + txpool.inspect, can be expensive
}

func (pool *dumbTxPool) Pending(enforceTips bool) (map[common.Address]types.Transactions, error) {
	// used in sync to send transactions
}

func (pool *dumbTxPool) Locals() []common.Address {

}
func (pool *dumbTxPool) AddLocal(tx *types.Transaction) error {

}
func (pool *dumbTxPool) AddRemotes(txs []*types.Transaction) []error {

}
func (pool *dumbTxPool) AddRemotesSync(txs []*types.Transaction) []error {

}
func (pool *dumbTxPool) ContentFrom(common.Address) (types.Transactions, types.Transactions) {

}
func (pool *dumbTxPool) Get(hash common.Hash) *types.Transaction {

}

func (pool *dumbTxPool) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return pool.scope.Track(pool.txFeed.Subscribe(ch))
}

func (pool *dumbTxPool) Stop() {
	pool.scope.Close()
	pool.chainHeadSub.Unsubscribe()
}
