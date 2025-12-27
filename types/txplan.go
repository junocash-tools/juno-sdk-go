package types

type TxPlanKind string

const (
	TxPlanKindWithdrawal TxPlanKind = "withdrawal"
	TxPlanKindSweep      TxPlanKind = "sweep"
	TxPlanKindRebalance  TxPlanKind = "rebalance"
)

type NoteRef struct {
	TxID        string `json:"txid"`
	ActionIndex uint32 `json:"action_index"`
}

type TxOutput struct {
	ToAddress      string `json:"to_address"`
	AmountZatoshis uint64 `json:"amount_zatoshis"`
	MemoHex        string `json:"memo_hex,omitempty"`
}

type TxPlan struct {
	Version     Version    `json:"version"`
	Kind        TxPlanKind `json:"kind"`
	WalletID    string     `json:"wallet_id"`
	Inputs      []NoteRef  `json:"inputs,omitempty"`
	Outputs     []TxOutput `json:"outputs"`
	FeeZatoshis uint64     `json:"fee_zatoshis"`
	Metadata    any        `json:"metadata,omitempty"`
}
