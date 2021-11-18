package starkex

import (
	"fmt"
	"testing"
)

const MOCK_PUBLIC_KEY = "3b865a18323b8d147a12c556bfb1d502516c325b1477a23ba6c77af31f020fd"
const MOCK_PRIVATE_KEY = "58c7d5a90b1776bde86ebac077e053ed85b0f7164f53b080304a531947f46e3"

func TestNewOrderSigner(t *testing.T) {
	param := OrderSignParam{
		NetworkId:  NETWORK_ID_ROPSTEN,
		Market:     "ETH-USD",
		Side:       "BUY",
		PositionId: 12345,
		HumanSize:  "145.0005",
		HumanPrice: "350.00067",
		LimitFee:   "0.125",
		ClientId:   "This is an ID that the client came up with to describe this order",
		Expiration: "2020-09-17T04:15:55.028Z",
	}
	sign, err := OrderSign(MOCK_PRIVATE_KEY, param)
	// 00cecbe513ecdbf782cd02b2a5efb03e58d5f63d15f2b840e9bc0029af04e8dd0090b822b16f50b2120e4ea9852b340f7936ff6069d02acca02f2ed03029ace5
	fmt.Println("sign,err", sign, err)
}

func TestNewWithdrawSigner(t *testing.T) {
	param := WithdrawSignParam{
		NetworkId:   NETWORK_ID_ROPSTEN,
		PositionId:  12345,
		HumanAmount: "49.478023",
		ClientId:    "This is an ID that the client came up with to describe this withdrawal",
		Expiration:  "2020-09-17T04:15:55.028Z",
	}
	sign, err := WithdrawSign(MOCK_PRIVATE_KEY, param)
	// 05e48c33f8205a5359c95f1bd7385c1c1f587e338a514298c07634c0b6c952ba0687d6980502a5d7fa84ef6fdc00104db22c43c7fb83e88ca84f19faa9ee3de1
	fmt.Println("sign,err", sign, err)
}

func TestNewTransferSigner(t *testing.T) {
	param := TransferSignParam{
		NetworkId:          NETWORK_ID_MAINNET,
		CreditAmount:       "1",
		DebitAmount:        "2",
		SenderPositionId:   12345,
		ReceiverPositionId: 67890,
		ReceiverPublicKey:  "04a9ecd28a67407c3cff8937f329ca24fd631b1d9ca2b9f2df47c7ebf72bf0b0",
		ReceiverAddress:    "0x1234567890123456789012345678901234567890",
		Expiration:         "2020-09-17T04:15:55.028Z",
		ClientId:           "This is an ID that the client came up with to describe this transfer",
	}
	sign, err := TransferSign(MOCK_PRIVATE_KEY, param)
	// 0278aeb361938d4c377950487bb770fc9464bf5352e19117c03243efad4e10a302bb3983e05676c7952caa4acdc1a83426d5c8cb8c56d7f6c477cfdafd37718a
	fmt.Println("sign,err", sign, err)
}

func TestGetTransferErc20Fact(t *testing.T) {
	recipient := "0x1234567890123456789012345678901234567890"
	tokenDecimals := 3
	humanAmount := "123.456"
	tokenAddress := "0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa"
	salt := "0x1234567890abcdef"
	// 34052387b5efb6132a42b244cff52a85a507ab319c414564d7a89207d4473672
	fact, err := GetTransferErc20Fact(recipient, tokenDecimals, humanAmount, tokenAddress, salt)
	fmt.Println("fact", fact, err)
}

func TestFactToCondition(t *testing.T) {
	fact := "cf9492ae0554c642b57f5d9cabee36fb512dd6b6629bdc51e60efb3118b8c2d8"
	addr := "0xe4a295420b58a4a7aa5c98920d6e8a0ef875b17a"
	condition := FactToCondition(addr, fact)
	// 0x4d794792504b063843afdf759534f5ed510a3ca52e7baba2e999e02349dd24
	fmt.Println("condition", condition.Text(16))
}

func TestNonceByClientId(t *testing.T) {
	nonce := NonceByClientId("")
	// 1723841828
	fmt.Println("nonce", nonce)
}
