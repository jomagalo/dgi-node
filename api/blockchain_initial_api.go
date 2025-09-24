package api

import (
	"github.com/jomagalo/dgi-node/state"
)

type BlockchainInitialApi struct {
	nodeState *state.NodeState
}

func NewBlockchainInitialApi(nodeState *state.NodeState) *BlockchainInitialApi {
	return &BlockchainInitialApi{nodeState}
}

func (api *BlockchainInitialApi) Syncing() Syncing {
	return Syncing{
		Syncing: true,
		Message: api.nodeState.Info(),
	}
}
