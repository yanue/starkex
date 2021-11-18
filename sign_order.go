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

type OrderSigner struct {
	param OrderSignParam
	msg   struct {
		OrderType               string   `json:"order_type"`
		AssetIdSynthetic        *big.Int `json:"asset_id_synthetic"`
		AssetIdCollateral       *big.Int `json:"asset_id_collateral"`
		AssetIdFee              *big.Int `json:"asset_id_fee"`
		QuantumAmountSynthetic  *big.Int `json:"quantum_amount_synthetic"`
		QuantumAmountCollateral *big.Int `json:"quantum_amount_collateral"`
		QuantumAmountFee        *big.Int `json:"quantum_amount_fee"`
		IsBuyingSynthetic       bool     `json:"is_buying_synthetic"`
		PositionId              *big.Int `json:"position_id"`
		Nonce                   *big.Int `json:"nonce"`
		ExpirationEpochHours    *big.Int `json:"expiration_epoch_hours"`
	}
}

func (s *OrderSigner) initMsg() error {
	currency := strings.Split(s.param.Market, "-")[0]                        // EOS-USD -> EOS
	assetIdSyn, ok := big.NewInt(0).SetString(SYNTHETIC_ID_MAP[currency], 0) // with prefix: 0x
	if !ok {
		return errors.New("invalid market: " + s.param.Market)
	}
	assetId := COLLATERAL_ASSET_ID_BY_NETWORK_ID[s.param.NetworkId] // asset id
	if assetId == nil {
		return errors.New(fmt.Sprintf("invalid network_id: %v", s.param.NetworkId))
	}
	exp, err := time.Parse("2006-01-02T15:04:05.000Z", s.param.Expiration)
	if err != nil {
		return err
	}
	resolutionC := decimal.NewFromInt(ASSET_RESOLUTION[currency])
	price, err := decimal.NewFromString(s.param.HumanPrice)
	if err != nil {
		return err
	}
	size, err := decimal.NewFromString(s.param.HumanSize)
	if err != nil {
		return err
	}
	var quantumsAmountSynthetic = decimal.NewFromFloat(0)
	isBuy := s.param.Side == "BUY"
	if isBuy {
		quantumsAmountSynthetic = size.Mul(price).Mul(resolutionUsdc).RoundUp(0)
	} else {
		quantumsAmountSynthetic = size.Mul(price).Mul(resolutionUsdc).RoundDown(0)
	}
	limitFeeRounded, err := decimal.NewFromString(s.param.LimitFee)
	if err != nil {
		return err
	}
	s.msg.OrderType = "LIMIT_ORDER_WITH_FEES"
	s.msg.AssetIdSynthetic = assetIdSyn
	s.msg.AssetIdCollateral = assetId
	s.msg.AssetIdFee = assetId
	s.msg.QuantumAmountSynthetic = size.Mul(resolutionC).BigInt()
	s.msg.QuantumAmountCollateral = quantumsAmountSynthetic.BigInt()
	s.msg.QuantumAmountFee = limitFeeRounded.Mul(quantumsAmountSynthetic).RoundUp(0).BigInt()
	s.msg.IsBuyingSynthetic = isBuy
	s.msg.PositionId = big.NewInt(s.param.PositionId)
	s.msg.Nonce = NonceByClientId(s.param.ClientId)
	s.msg.ExpirationEpochHours = big.NewInt(int64(math.Ceil(float64(exp.Unix())/float64(ONE_HOUR_IN_SECONDS))) + ORDER_SIGNATURE_EXPIRATION_BUFFER_HOURS)
	return nil
}

func (s *OrderSigner) getHash() (string, error) {
	var assetIdSell, assetIdBuy, quantumsAmountSell, quantumsAmountBuy *big.Int
	if s.msg.IsBuyingSynthetic {
		assetIdSell = s.msg.AssetIdCollateral
		assetIdBuy = s.msg.AssetIdSynthetic
		quantumsAmountSell = s.msg.QuantumAmountCollateral
		quantumsAmountBuy = s.msg.QuantumAmountSynthetic
	} else {
		assetIdSell = s.msg.AssetIdSynthetic
		assetIdBuy = s.msg.AssetIdCollateral
		quantumsAmountSell = s.msg.QuantumAmountSynthetic
		quantumsAmountBuy = s.msg.QuantumAmountCollateral
	}
	fee := s.msg.QuantumAmountFee
	nonce := s.msg.Nonce
	// part1
	part1 := big.NewInt(0).Set(quantumsAmountSell)
	part1.Lsh(part1, ORDER_FIELD_BIT_LENGTHS["quantums_amount"])
	part1.Add(part1, quantumsAmountBuy)
	part1.Lsh(part1, ORDER_FIELD_BIT_LENGTHS["quantums_amount"])
	part1.Add(part1, fee)
	part1.Lsh(part1, ORDER_FIELD_BIT_LENGTHS["nonce"])
	part1.Add(part1, nonce)
	// part2
	part2 := big.NewInt(ORDER_PREFIX)
	for i := 0; i < 3; i++ {
		part2.Lsh(part2, ORDER_FIELD_BIT_LENGTHS["position_id"])
		part2.Add(part2, s.msg.PositionId)
	}
	part2.Lsh(part2, ORDER_FIELD_BIT_LENGTHS["expiration_epoch_hours"])
	part2.Add(part2, s.msg.ExpirationEpochHours)
	part2.Lsh(part2, ORDER_PADDING_BITS)
	// pedersen hash
	assetHash := getHash(getHash(assetIdSell.String(), assetIdBuy.String()), s.msg.AssetIdFee.String())
	part1Hash := getHash(assetHash, part1.String())
	part2Hash := getHash(part1Hash, part2.String())
	return part2Hash, nil
}
