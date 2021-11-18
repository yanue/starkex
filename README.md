## stark key authentication library, signature generator for dydx exchange implemented in pure golang

for the following operations:

- Place an order
- Withdraw funds

> link : https://docs.dydx.exchange/#authentication

### Installation

```
go get github.com/yanue/starkex
```

### demos

#### order sign demo

```
    const MOCK_PRIVATE_KEY = "58c7d5a90b1776bde86ebac077e053ed85b0f7164f53b080304a531947f46e3"
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
```

#### withdraw sign demo

```
    const MOCK_PRIVATE_KEY = "58c7d5a90b1776bde86ebac077e053ed85b0f7164f53b080304a531947f46e3"
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
```

#### transfer sign demo (fast_withdraw)

```
    const MOCK_PRIVATE_KEY = "58c7d5a90b1776bde86ebac077e053ed85b0f7164f53b080304a531947f46e3"
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
```

### inspired by

> https://github.com/dydxprotocol/dydx-v3-python
> https://github.com/starkware-libs/starkex-resources