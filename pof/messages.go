package pof

import (
	"github.com/fletaio/fleta_testnet/common"
	"github.com/fletaio/fleta_testnet/common/hash"
	"github.com/fletaio/fleta_testnet/core/types"
)

// RoundVote is a message for a round vote
type RoundVote struct {
	ChainID              uint8
	LastHash             hash.Hash256
	TargetHeight         uint32
	TimeoutCount         uint32
	Formulator           common.Address
	FormulatorPublicHash common.PublicHash
	Timestamp            uint64
	IsReply              bool
}

// RoundVoteMessage is a message for a round vote
type RoundVoteMessage struct {
	RoundVote *RoundVote
	Signature common.Signature
}

// RoundVoteAck is a message for a round vote ack
type RoundVoteAck struct {
	ChainID              uint8
	LastHash             hash.Hash256
	TargetHeight         uint32
	TimeoutCount         uint32
	Formulator           common.Address
	FormulatorPublicHash common.PublicHash
	PublicHash           common.PublicHash
	Timestamp            uint64
	IsReply              bool
}

// RoundVoteAckMessage is a message for a round vote
type RoundVoteAckMessage struct {
	RoundVoteAck *RoundVoteAck
	Signature    common.Signature
}

// RoundSetupMessage is a message for a round vote
type RoundSetupMessage struct {
	MinRoundVoteAck *RoundVoteAck
	Signature       common.Signature
}

// BlockReqMessage is a message for a block request
type BlockReqMessage struct {
	PrevHash             hash.Hash256
	TargetHeight         uint32
	TimeoutCount         uint32
	Formulator           common.Address
	FormulatorPublicHash common.PublicHash
}

// BlockGenMessage is a message for a block generation
type BlockGenMessage struct {
	Block              *types.Block
	GeneratorSignature common.Signature
	IsReply            bool
}

// BlockVote is message for a block vote
type BlockVote struct {
	TargetHeight       uint32
	Header             *types.Header
	GeneratorSignature common.Signature
	ObserverSignature  common.Signature
	IsReply            bool
}

// BlockVoteMessage is a message for a round vote
type BlockVoteMessage struct {
	BlockVote *BlockVote
	Signature common.Signature
}

// BlockObSignMessage is a message for a block observer signatures
type BlockObSignMessage struct {
	TargetHeight       uint32
	BlockSign          *types.BlockSign
	ObserverSignatures []common.Signature
}

// BlockGenRequest is a message to request block gen
type BlockGenRequest struct {
	ChainID              uint8
	LastHash             hash.Hash256
	TargetHeight         uint32
	TimeoutCount         uint32
	Formulator           common.Address
	FormulatorPublicHash common.PublicHash
	PublicHash           common.PublicHash
	Timestamp            uint64
}

// BlockGenRequestMessage is a message to request block gen
type BlockGenRequestMessage struct {
	BlockGenRequest *BlockGenRequest
	Signature       common.Signature
}
