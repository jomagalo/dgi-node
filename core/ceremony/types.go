package ceremony

import (
	"github.com/jomagalo/dgi-node/blockchain/types"
	"github.com/jomagalo/dgi-node/common"
	"github.com/shopspring/decimal"
)

type candidate struct {
	Address  common.Address
	PubKey   []byte
	IsAuthor bool
}

type candidatesOfShard struct {
	candidates     []*candidate
	flips          [][]byte
	flipsPerAuthor map[int][][]byte
	flipAuthorMap  map[string]common.Address

	nonCandidates []common.Address

	shortFlipsPerCandidate [][]int
	longFlipsPerCandidate  [][]int

	shortFlipsToSolve map[common.Address][][]byte
	longFlipsToSolve  map[common.Address][][]byte
}

type FlipStatus byte

const (
	NotQualified    FlipStatus = 0
	Qualified       FlipStatus = 1
	WeaklyQualified FlipStatus = 2
	QualifiedByNone FlipStatus = 3
)

type FlipQualification struct {
	status     FlipStatus
	answer     types.Answer
	grade      types.Grade
	gradeScore decimal.Decimal
}
