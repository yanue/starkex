package starkex

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"math/big"
	"strings"
	"time"
)

type TransferSigner struct {
	param TransferSignParam
	msg   struct {
		SenderPositionId     *big.Int `json:"sender_position_id"`
		ReceiverPositionId   *big.Int `json:"receiver_position_id"`
		ReceiverPublicKey    *big.Int `json:"receiver_public_key"`
		Condition            *big.Int `json:"condition"`
		QuantumsAmount       *big.Int `json:"quantums_amount"`
		Nonce                *big.Int `json:"nonce"`
		ExpirationEpochHours *big.Int `json:"expiration_epoch_hours"`
	}
}

func (s *TransferSigner) initMsg() error {
	exp, err := time.Parse("2006-01-02T15:04:05.000Z", s.param.Expiration)
	if err != nil {
		return err
	}
	QuantumAmount, err := decimal.NewFromString(s.param.DebitAmount)
	if err != nil {
		return err
	}
	receiverKey, ok := big.NewInt(0).SetString(strings.TrimPrefix(s.param.ReceiverPublicKey, "0x"), 16)
	if !ok {
		return errors.New(fmt.Sprintf("invalid receiver_public_key: %v", s.param.ReceiverPublicKey))
	}
	fact, err := s.getFact()
	if err != nil {
		return err
	}
	factRegistryAddress := FACT_REGISTRY_CONTRACT[s.param.NetworkId]
	condition := FactToCondition(factRegistryAddress, fact)
	// set msg
	s.msg.QuantumsAmount = QuantumAmount.Mul(resolutionUsdc).BigInt()
	s.msg.Condition = condition
	s.msg.SenderPositionId = big.NewInt(s.param.SenderPositionId)
	s.msg.ReceiverPositionId = big.NewInt(s.param.ReceiverPositionId)
	s.msg.Nonce = NonceByClientId(s.param.ClientId)
	s.msg.ReceiverPublicKey = receiverKey
	s.msg.ExpirationEpochHours = big.NewInt(int64(math.Ceil(float64(exp.Unix()) / float64(ONE_HOUR_IN_SECONDS))))
	return nil
}

func (s *TransferSigner) getFact() (string, error) {
	// generate
	salt := NonceByClientId(s.param.ClientId).String()
	tokenAddr := TOKEN_CONTRACTS[COLLATERAL_ASSET][s.param.NetworkId]
	fact, err := GetTransferErc20Fact(s.param.ReceiverAddress, COLLATERAL_TOKEN_DECIMALS, s.param.CreditAmount, tokenAddr, salt)
	return fact, err
}

func (s *TransferSigner) getHash() (string, error) {
	// net
	net := COLLATERAL_ASSET_ID_BY_NETWORK_ID[s.param.NetworkId]
	assetHash := getHash(net.String(), big.NewInt(CONDITIONAL_TRANSFER_FEE_ASSET_ID).String())
	keyHash := getHash(assetHash, s.msg.ReceiverPublicKey.String())
	// part 1
	part1 := getHash(keyHash, s.msg.Condition.String())
	// part 2
	part2 := big.NewInt(0).Set(s.msg.SenderPositionId)
	part2.Lsh(part2, CONDITIONAL_TRANSFER_FIELD_BIT_LENGTHS["position_id"])
	part2.Add(part2, s.msg.ReceiverPositionId)
	part2.Lsh(part2, CONDITIONAL_TRANSFER_FIELD_BIT_LENGTHS["position_id"])
	part2.Add(part2, s.msg.SenderPositionId)
	part2.Lsh(part2, CONDITIONAL_TRANSFER_FIELD_BIT_LENGTHS["nonce"])
	part2.Add(part2, s.msg.Nonce)
	// part 3
	part3 := big.NewInt(CONDITIONAL_TRANSFER_PREFIX)
	part3.Lsh(part3, CONDITIONAL_TRANSFER_FIELD_BIT_LENGTHS["quantums_amount"])
	part3.Add(part3, s.msg.QuantumsAmount)
	part3.Lsh(part3, CONDITIONAL_TRANSFER_FIELD_BIT_LENGTHS["quantums_amount"])
	part3.Add(part3, big.NewInt(CONDITIONAL_TRANSFER_MAX_AMOUNT_FEE))
	part3.Lsh(part3, CONDITIONAL_TRANSFER_FIELD_BIT_LENGTHS["expiration_epoch_hours"])
	part3.Add(part3, s.msg.ExpirationEpochHours)
	part3.Lsh(part3, CONDITIONAL_TRANSFER_PADDING_BITS)
	// pedersen hash
	hash1 := getHash(part1, part2.String())
	hash2 := getHash(hash1, part3.String())
	return hash2, nil
}
