package starkex

type OrderSignParam struct {
	NetworkId  int    `json:"network_id"` // 1 MAINNET 3 ROPSTEN
	PositionId int64  `json:"position_id"`
	Market     string `json:"market"`
	Side       string `json:"side"`
	HumanSize  string `json:"human_size"`
	HumanPrice string `json:"human_price"`
	LimitFee   string `json:"limit_fee"`
	ClientId   string `json:"clientId"`
	Expiration string `json:"expiration"` // 2006-01-02T15:04:05.000Z
}

type WithdrawSignParam struct {
	NetworkId   int    `json:"network_id"` // 1 MAINNET 3 ROPSTEN
	PositionId  int64  `json:"position_id"`
	HumanAmount string `json:"human_amount"`
	ClientId    string `json:"clientId"`
	Expiration  string `json:"expiration"` // 2006-01-02T15:04:05.000Z
}

type TransferSignParam struct {
	NetworkId          int    `json:"network_id"` // 1 MAINNET 3 ROPSTEN
	SenderPositionId   int64  `json:"sender_position_id"`
	ReceiverPositionId int64  `json:"receiver_position_id"`
	ReceiverPublicKey  string `json:"receiver_public_key"`
	ReceiverAddress    string `json:"receiver_address"`
	CreditAmount       string `json:"credit_amount"`
	DebitAmount        string `json:"debit_amount"`
	Expiration         string `json:"expiration"`
	ClientId           string `json:"client_id"`
}
