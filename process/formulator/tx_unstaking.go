package formulator

import (
	"bytes"
	"encoding/json"

	"github.com/fletaio/fleta_testnet/common"
	"github.com/fletaio/fleta_testnet/common/amount"
	"github.com/fletaio/fleta_testnet/core/types"
	"github.com/fletaio/fleta_testnet/encoding"
)

// Unstaking is used to ustake coin from the hyper formulator
type Unstaking struct {
	Timestamp_      uint64
	Seq_            uint64
	From_           common.Address
	HyperFormulator common.Address
	Amount          *amount.Amount
}

// Timestamp returns the timestamp of the transaction
func (tx *Unstaking) Timestamp() uint64 {
	return tx.Timestamp_
}

// Seq returns the sequence of the transaction
func (tx *Unstaking) Seq() uint64 {
	return tx.Seq_
}

// From returns the from address of the transaction
func (tx *Unstaking) From() common.Address {
	return tx.From_
}

// Fee returns the fee of the transaction
func (tx *Unstaking) Fee(loader types.LoaderWrapper) *amount.Amount {
	return amount.COIN.DivC(10)
}

// Validate validates signatures of the transaction
func (tx *Unstaking) Validate(p types.Process, loader types.LoaderWrapper, signers []common.PublicHash) error {
	sp := p.(*Formulator)

	if tx.Amount.Less(amount.COIN) {
		return ErrInvalidStakingAmount
	}

	if tx.Seq() <= loader.Seq(tx.From()) {
		return types.ErrInvalidSequence
	}

	acc, err := loader.Account(tx.HyperFormulator)
	if err != nil {
		return err
	}
	frAcc, is := acc.(*FormulatorAccount)
	if !is {
		return types.ErrInvalidAccountType
	}
	if frAcc.FormulatorType != HyperFormulatorType {
		return types.ErrInvalidAccountType
	}

	fromAcc, err := loader.Account(tx.From())
	if err != nil {
		return err
	}
	if err := fromAcc.Validate(loader, signers); err != nil {
		return err
	}

	fromStakingAmount := sp.GetStakingAmount(loader, tx.HyperFormulator, tx.From())
	if fromStakingAmount.Less(tx.Amount) {
		return ErrInsufficientStakingAmount
	}
	if frAcc.StakingAmount.Less(tx.Amount) {
		return ErrInsufficientStakingAmount
	}

	if err := sp.vault.CheckFeePayable(loader, tx); err != nil {
		return err
	}
	return nil
}

// Execute updates the context by the transaction
func (tx *Unstaking) Execute(p types.Process, ctw *types.ContextWrapper, index uint16) error {
	sp := p.(*Formulator)

	return sp.vault.WithFee(ctw, tx, func() error {
		acc, err := ctw.Account(tx.HyperFormulator)
		if err != nil {
			return err
		}
		frAcc := acc.(*FormulatorAccount)

		if err := sp.subStakingAmount(ctw, tx.HyperFormulator, tx.From(), tx.Amount); err != nil {
			return err
		}
		frAcc.StakingAmount = frAcc.StakingAmount.Sub(tx.Amount)

		policy := &HyperPolicy{}
		if err := encoding.Unmarshal(ctw.ProcessData(tagHyperPolicy), &policy); err != nil {
			return err
		}
		if err := sp.addUnstakingAmount(ctw, tx.HyperFormulator, tx.From(), ctw.TargetHeight()+policy.StakingUnlockRequiredBlocks, tx.Amount); err != nil {
			return err
		}
		return nil
	})
}

// MarshalJSON is a marshaler function
func (tx *Unstaking) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString(`{`)
	buffer.WriteString(`"timestamp":`)
	if bs, err := json.Marshal(tx.Timestamp_); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`,`)
	buffer.WriteString(`"seq":`)
	if bs, err := json.Marshal(tx.Seq_); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`,`)
	buffer.WriteString(`"from":`)
	if bs, err := tx.From_.MarshalJSON(); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`,`)
	buffer.WriteString(`"hyper_formulator":`)
	if bs, err := tx.HyperFormulator.MarshalJSON(); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`,`)
	buffer.WriteString(`"amount":`)
	if bs, err := tx.Amount.MarshalJSON(); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`}`)
	return buffer.Bytes(), nil
}
