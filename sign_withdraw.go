package starkex

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"math/big"
	"time"
)

// Sign for withdraw

type WithdrawSigner struct {
	param WithdrawSignParam
	msg   struct {
		PositionId           *big.Int `json:"position_id"`
		QuantumAmount        *big.Int `json:"quantum_amount"`
		Nonce                *big.Int `json:"nonce"`
		ExpirationEpochHours *big.Int `json:"expiration_epoch_hours"`
	}
}

func (s *WithdrawSigner) initMsg() error {
	exp, err := time.Parse("2006-01-02T15:04:05.000Z", s.param.Expiration)
	if err != nil {
		return err
	}
	QuantumAmount, err := decimal.NewFromString(s.param.HumanAmount)
	if err != nil {
		return err
	}
	s.msg.QuantumAmount = QuantumAmount.Mul(resolutionUsdc).BigInt()
	s.msg.PositionId = big.NewInt(s.param.PositionId)
	s.msg.Nonce = NonceByClientId(s.param.ClientId)
	s.msg.ExpirationEpochHours = big.NewInt(int64(math.Ceil(float64(exp.Unix()) / float64(ONE_HOUR_IN_SECONDS))))
	return nil
}

func (s *WithdrawSigner) getHash() (string, error) {
	net := COLLATERAL_ASSET_ID_BY_NETWORK_ID[s.param.NetworkId]
	if net == nil {
		return "", errors.New(fmt.Sprintf("invalid network_id: %v", s.param.NetworkId))
	}
	// packed
	packed := big.NewInt(WITHDRAWAL_PREFIX)
	packed.Lsh(packed, WITHDRAWAL_FIELD_BIT_LENGTHS["position_id"])
	packed.Add(packed, s.msg.PositionId)
	packed.Lsh(packed, WITHDRAWAL_FIELD_BIT_LENGTHS["nonce"])
	packed.Add(packed, s.msg.Nonce)
	packed.Lsh(packed, WITHDRAWAL_FIELD_BIT_LENGTHS["quantums_amount"])
	packed.Add(packed, s.msg.QuantumAmount)
	packed.Lsh(packed, WITHDRAWAL_FIELD_BIT_LENGTHS["expiration_epoch_hours"])
	packed.Add(packed, s.msg.ExpirationEpochHours)
	packed.Lsh(packed, WITHDRAWAL_PADDING_BITS)
	// pedersen hash
	hash := getHash(net.String(), packed.String())
	return hash, nil
}
