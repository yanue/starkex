package starkex

import (
	"encoding/json"
	"math/big"
)

/**
# Starkware crypto functions implemented in Golang.
#
# python source:
# https://github.com/starkware-libs/starkex-resources/blob/0f08e6c55ad88c93499f71f2af4a2e7ae0185cdf/crypto/starkware/crypto/signature/signature.py
*/

type PedersenCfg struct {
	Comment        string        `json:"_comment"`
	FieldPrime     *big.Int      `json:"FIELD_PRIME"`
	FieldGen       int           `json:"FIELD_GEN"`
	EcOrder        *big.Int      `json:"EC_ORDER"`
	ALPHA          int           `json:"ALPHA"`
	BETA           *big.Int      `json:"BETA"`
	ConstantPoints [][2]*big.Int `json:"CONSTANT_POINTS"`
}

var pedersenCfg PedersenCfg

var EC_ORDER = new(big.Int)
var FIELD_PRIME = new(big.Int)

func init() {
	_ = json.Unmarshal([]byte(pedersenParams), &pedersenCfg)
	EC_ORDER = pedersenCfg.EcOrder
	FIELD_PRIME = pedersenCfg.FieldPrime
}

func PedersenHash(str ...string) string {
	NElementBitsHash := FIELD_PRIME.BitLen()
	point := pedersenCfg.ConstantPoints[0]
	for i, s := range str {
		x, _ := big.NewInt(0).SetString(s, 10)
		pointList := pedersenCfg.ConstantPoints[2+i*NElementBitsHash : 2+(i+1)*NElementBitsHash]
		n := big.NewInt(0)
		for _, pt := range pointList {
			n.And(x, big.NewInt(1))
			if n.Cmp(big.NewInt(0)) > 0 {
				point = eccAdd(point, pt, FIELD_PRIME)
			}
			x = x.Rsh(x, 1)
		}
	}
	return point[0].String()
}
