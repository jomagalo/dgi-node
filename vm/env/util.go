package env

import (
	"github.com/jomagalo/dgi-node/blockchain/types"
	"github.com/jomagalo/dgi-node/common"
	"github.com/jomagalo/dgi-node/crypto"
)

func ComputeContractAddr(tx *types.Transaction, from common.Address) common.Address {
	hash := crypto.Hash(append(append(from.Bytes(), common.ToBytes(tx.Epoch)...), common.ToBytes(tx.AccountNonce)...))
	var result common.Address
	result.SetBytes(hash[:])
	return result
}