package blockchain

import (
	"github.com/jomagalo/dgi-node/blockchain/attachments"
	"github.com/jomagalo/dgi-node/blockchain/types"
	"github.com/jomagalo/dgi-node/common"
	"github.com/jomagalo/dgi-node/common/eventbus"
	"github.com/jomagalo/dgi-node/config"
	"github.com/jomagalo/dgi-node/database"
	"github.com/jomagalo/dgi-node/events"
	"github.com/jomagalo/dgi-node/keystore"
	dbm "github.com/tendermint/tm-db"
	"sync"
)

type indexer struct {
	coinbase common.Address
	repo     *database.Repo
	bus      eventbus.Bus
	keystore *keystore.KeyStore
	cfg      *config.Config
	mutex    sync.Mutex
}

func newBlockchainIndexer(db dbm.DB, bus eventbus.Bus, cfg *config.Config, keystore *keystore.KeyStore) *indexer {
	return &indexer{
		repo:     database.NewRepo(db),
		bus:      bus,
		keystore: keystore,
		cfg:      cfg,
	}
}

func (i *indexer) initialize(coinbase common.Address) {
	i.coinbase = coinbase
}

func (i *indexer) HandleBlockTransactions(header *types.Header, txs []*types.Transaction) {

	i.repo.DeleteOutdatedBurntCoins(header.Height(), i.cfg.Blockchain.BurnTxRange)

	accounts := i.keystore.Accounts()
	accountsMap := make(map[common.Address]struct{})
	for _, item := range accounts {
		accountsMap[item.Address] = struct{}{}
	}
	accountsMap[i.coinbase] = struct{}{}

	for _, tx := range txs {
		sender, _ := types.Sender(tx)
		i.handleOwnTx(header, sender, tx, accountsMap)
		i.handleBurnTx(header.Height(), sender, tx)
		i.handleOwnDeleteFlipTx(sender, tx)
	}
}

func (i *indexer) handleOwnTx(header *types.Header, sender common.Address, tx *types.Transaction, accountsMap map[common.Address]struct{}) {
	if _, ok := accountsMap[sender]; ok {
		i.repo.SaveTx(sender, header.Hash(), header.Time(), header.FeePerGas(), tx)
	}
	if tx.To != nil {
		to := *tx.To
		if sender == to {
			return
		}
		if _, ok := accountsMap[to]; ok {
			i.repo.SaveTx(to, header.Hash(), header.Time(), header.FeePerGas(), tx)
		}
	}
}

func (i *indexer) handleBurnTx(height uint64, sender common.Address, tx *types.Transaction) {
	if i.cfg.Consensus.EnableUpgrade10 {
		return
	}
	if tx.Type != types.BurnTx {
		return
	}
	attachment := attachments.ParseBurnAttachment(tx)
	if attachment == nil {
		return
	}
	i.repo.SaveBurntCoins(height, tx.Hash(), sender, attachment.Key, tx.AmountOrZero())
}

func (i *indexer) handleOwnDeleteFlipTx(sender common.Address, tx *types.Transaction) {
	if tx.Type != types.DeleteFlipTx || sender != i.coinbase {
		return
	}
	attachment := attachments.ParseDeleteFlipAttachment(tx)
	if attachment == nil {
		return
	}
	i.bus.Publish(&events.DeleteFlipEvent{
		FlipCid: attachment.Cid,
	})
}
