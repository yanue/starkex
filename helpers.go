package starkex

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/huandu/xstrings"
	"github.com/miguelmota/go-solidity-sha3"
	"github.com/shopspring/decimal"
	"log"
	"math/big"
	"strings"
)

func ToJsonString(input interface{}) string {
	js, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		log.Println("ToJsonString error:", err.Error())
	}
	return string(js)
}

func getHash(str1, str2 string) string {
	return PedersenHash(str1, str2)
}

// NonceByClientId generate nonce by clientId
func NonceByClientId(clientId string) *big.Int {
	h := sha256.New()
	h.Write([]byte(clientId))

	a := new(big.Int)
	a.SetBytes(h.Sum(nil))
	res := a.Mod(a, big.NewInt(NONCE_UPPER_BOUND_EXCLUSIVE))
	return res
}

// SerializeSignature Convert a Sign from an r, s pair to a 32-byte hex string.
func SerializeSignature(r, s *big.Int) string {
	return IntToHex32(r) + IntToHex32(s)
}

// IntToHex32 Normalize to a 32-byte hex string without 0x prefix.
func IntToHex32(x *big.Int) string {
	str := x.Text(16)
	return xstrings.RightJustify(str, 64, "0")
}

// FactToCondition Generate the condition, signed as part of a conditional transfer.
func FactToCondition(factRegistryAddress string, fact string) *big.Int {
	data := strings.TrimPrefix(factRegistryAddress, "0x") + fact
	hexBytes, _ := hex.DecodeString(data)
	// int(Web3.keccak(data).hex(), 16) & BIT_MASK_250
	hash := crypto.Keccak256Hash(hexBytes)
	fst := hash.Big()
	fst.And(fst, BIT_MASK_250)
	return fst
}

// GetTransferErc20Fact get erc20 fact
// tokenDecimals is COLLATERAL_TOKEN_DECIMALS
func GetTransferErc20Fact(recipient string, tokenDecimals int, humanAmount, tokenAddress, salt string) (string, error) {
	fmt.Println("GetTransferErc20Fact", recipient, tokenDecimals, humanAmount, tokenAddress, salt)
	// token_amount = float(human_amount) * (10 ** token_decimals)
	amount, err := decimal.NewFromString(humanAmount)
	if err != nil {
		return "", err
	}
	saltInt, ok := big.NewInt(0).SetString(salt, 0) // with prefix: 0x
	if !ok {
		return "", errors.New(fmt.Sprintf("invalid salt: %v,can not parse to big.Int", salt))
	}
	tokenAmount := amount.Mul(decimal.New(10, int32(tokenDecimals-1)))
	fact := solsha3.SoliditySHA3(
		// types
		[]string{"address", "uint256", "address", "uint256"},
		// values
		[]interface{}{recipient, tokenAmount.String(), tokenAddress, saltInt.String()},
	)
	return hex.EncodeToString(fact), nil
}

func GenerateKRfc6979(msgHash, priKey *big.Int, seed int) *big.Int {
	msgHash = big.NewInt(0).Set(msgHash) // copy
	bitMod := msgHash.BitLen() % 8
	if bitMod <= 4 && bitMod >= 1 && msgHash.BitLen() > 248 {
		msgHash.Mul(msgHash, big.NewInt(16))
	}
	var extra []byte
	if seed > 0 {
		buf := new(bytes.Buffer)
		var data interface{}
		if seed < 256 {
			data = uint8(seed)
		} else if seed < 65536 {
			data = uint16(seed)
		} else if seed < 4294967296 {
			data = uint32(seed)
		} else {
			data = uint64(seed)
		}
		_ = binary.Write(buf, binary.BigEndian, data)
		extra = buf.Bytes()
	}
	return generateSecret(EC_ORDER, priKey, sha256.New, msgHash.Bytes(), extra)
}
