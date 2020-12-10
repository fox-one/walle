package order

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/vmihailenco/msgpack/v5"
)

//go:generate stringer -type TransactionType -trimprefix TransactionType

type TransactionType int

const (
	_ TransactionType = iota
	TransactionTypeAdd
	TransactionTypeRemove
	TransactionTypeSwap
)

func ParseTransactionType(t string) TransactionType {
	for idx := 0; idx < len(_TransactionType_index)-1; idx++ {
		l, r := _TransactionType_index[idx], _TransactionType_index[idx+1]
		if typ := _TransactionType_name[l:r]; strings.EqualFold(typ, t) {
			return TransactionType(idx + 1)
		}
	}

	return 0
}

type TransactionAction struct {
	// Transaction type add remove swap
	Type string `json:"t,omitempty" msgpack:"t,omitempty"`
	// deposit
	Deposit string `json:"d,omitempty" msgpack:"d,omitempty"`
	// withdraw
	Pairs         []string `json:"p,omitempty" msgpack:"p,omitempty"`
	RemovePercent int64    `json:"l,omitempty" msgpack:"l,omitempty"`
	// Swap
	AssetID string `json:"a,omitempty" msgpack:"a,omitempty"`
	Routes  string `json:"r,omitempty" msgpack:"r,omitempty"`
	Minimum string `json:"m,omitempty" msgpack:"m,omitempty"`
}

func ValidateTransactionAction(action TransactionAction) error {
	typ := ParseTransactionType(action.Type)
	switch typ {
	case TransactionTypeAdd:
		if !govalidator.IsUUID(action.Deposit) {
			return errors.New("deposit id must be uuid")
		}

		if !govalidator.IsUUID(action.AssetID) {
			return errors.New("asset id must be uuid")
		}
	case TransactionTypeRemove:
		if len(action.Pairs) > 0 {
			if len(action.Pairs) != 2 {
				return errors.New("invalid pair assets")
			}

			base, quote := action.Pairs[0], action.Pairs[1]
			if !govalidator.IsUUID(base) {
				return errors.New("base asset id must be uuid")
			}

			if !govalidator.IsUUID(quote) {
				return errors.New("quote asset id must be uuid")
			}

			if base == quote {
				return errors.New("base asset id is same with quote asset id")
			}

			if action.RemovePercent < 1 || action.RemovePercent > 100 {
				return errors.New("remove percent must in [1,100]")
			}
		}
	case TransactionTypeSwap:
		if !govalidator.IsUUID(action.AssetID) {
			return errors.New("asset id must be uuid")
		}
	default:
		return fmt.Errorf("invalid transaction type %s", typ)
	}

	return nil
}

func DecodeTransactionAction(memo string) (*TransactionAction, error) {
	var action TransactionAction
	if err := DecodeAction(&action, memo); err != nil {
		return nil, err
	}

	if err := ValidateTransactionAction(action); err != nil {
		return nil, err
	}

	return &action, nil
}

func DecodeAction(v interface{}, memo string) error {
	b, err := base64.StdEncoding.DecodeString(memo)

	if err != nil {
		b, err = base64.URLEncoding.DecodeString(memo)
	}

	if err != nil {
		b = []byte(memo)
	}

	if err := msgpack.Unmarshal(b, v); err == nil {
		return nil
	}

	return json.Unmarshal(b, v)
}
