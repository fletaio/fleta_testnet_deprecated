package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strings"

	"github.com/fletaio/fleta_testnet/common"
	"github.com/petar/GoLLRB/llrb"
)

var (
	ninf = nInf{}
	pinf = pInf{}
)

type nInf struct{}

func (nInf) Less(llrb.Item) bool {
	return true
}

type pInf struct{}

func (pInf) Less(llrb.Item) bool {
	return false
}

func cmpAddressASC(a interface{}, b interface{}) bool {
	ai := a.(common.Address)
	bi := b.(common.Address)
	return bytes.Compare(ai[:], bi[:]) < 0
}

func cmpStringASC(a interface{}, b interface{}) bool {
	ai := a.(string)
	bi := b.(string)
	return strings.Compare(ai[:], bi[:]) < 0
}

func cmpUint64ASC(a interface{}, b interface{}) bool {
	ai := a.(uint64)
	bi := b.(uint64)
	return ai < bi
}

// IsAllowedAccountName returns it is allowed account name or not
func IsAllowedAccountName(Name string) bool {
	if len(Name) < 8 || len(Name) > 40 {
		return false
	}
	if _, err := common.ParseAddress(Name); err == nil {
		return false
	}
	for i := 0; i < len(Name); i++ {
		c := Name[i]
		if (c < '0' || '9' < c) && (c < 'a' || 'z' < c) && (c < 'A' || 'Z' < c) && c != '.' && c != '-' && c != '_' && c != '@' {
			return false
		}
	}
	return true
}

// UnmarshalID returns the block height, the transaction index in the block, the output index in the transaction
func UnmarshalID(id uint64) (uint32, uint16, uint16) {
	return uint32(id >> 32), uint16(id >> 16), uint16(id)
}

// MarshalID returns the packed id
func MarshalID(height uint32, index uint16, n uint16) uint64 {
	return uint64(height)<<32 | uint64(index)<<16 | uint64(n)
}

// TransactionID returns the id of the transaction
func TransactionID(Height uint32, Index uint16) string {
	bs := make([]byte, 6)
	binary.BigEndian.PutUint32(bs, Height)
	binary.BigEndian.PutUint16(bs[4:], Index)
	return hex.EncodeToString(bs)
}

// ParseTransactionID returns the id of the transaction
func ParseTransactionID(TXID string) (uint32, uint16, error) {
	if len(TXID) != 12 {
		return 0, 0, ErrInvalidTransactionIDFormat
	}
	bs, err := hex.DecodeString(TXID)
	if err != nil {
		return 0, 0, err
	}
	Height := binary.BigEndian.Uint32(bs)
	Index := binary.BigEndian.Uint16(bs[4:])
	return Height, Index, nil
}
