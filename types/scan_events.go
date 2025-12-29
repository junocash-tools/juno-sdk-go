package types

import "encoding/json"

// WalletEventKind is the "kind" discriminator used by juno-scan wallet events and broker envelopes.
type WalletEventKind string

const (
	WalletEventKindDepositEvent       WalletEventKind = "DepositEvent"
	WalletEventKindDepositConfirmed   WalletEventKind = "DepositConfirmed"
	WalletEventKindDepositOrphaned    WalletEventKind = "DepositOrphaned"
	WalletEventKindDepositUnconfirmed WalletEventKind = "DepositUnconfirmed"

	WalletEventKindSpendEvent       WalletEventKind = "SpendEvent"
	WalletEventKindSpendConfirmed   WalletEventKind = "SpendConfirmed"
	WalletEventKindSpendOrphaned    WalletEventKind = "SpendOrphaned"
	WalletEventKindSpendUnconfirmed WalletEventKind = "SpendUnconfirmed"
)

// BrokerEnvelope is the top-level message format emitted by juno-scan publisher adapters
// (Kafka/NATS/RabbitMQ).
type BrokerEnvelope struct {
	Version  Version         `json:"version"`
	Kind     WalletEventKind `json:"kind"`
	WalletID string          `json:"wallet_id"`
	Height   int64           `json:"height"`
	Payload  json.RawMessage `json:"payload"`
}

type DepositEventPayload struct {
	DepositEvent
	RecipientAddress string `json:"recipient_address,omitempty"`
	NoteNullifier    string `json:"note_nullifier,omitempty"`
}

type DepositConfirmedPayload struct {
	DepositEventPayload
	ConfirmedHeight       int64 `json:"confirmed_height"`
	RequiredConfirmations int64 `json:"required_confirmations"`
}

type DepositOrphanedPayload struct {
	DepositEventPayload
	OrphanedAtHeight int64 `json:"orphaned_at_height"`
}

type DepositUnconfirmedPayload struct {
	DepositEventPayload
	RollbackHeight          int64 `json:"rollback_height"`
	RequiredConfirmations   int64 `json:"required_confirmations,omitempty"`
	PreviousConfirmedHeight int64 `json:"previous_confirmed_height"`
}

type SpendEventPayload struct {
	Version          Version  `json:"version"`
	WalletID         string   `json:"wallet_id"`
	DiversifierIndex uint32   `json:"diversifier_index,omitempty"`
	TxID             string   `json:"txid"`
	Height           int64    `json:"height"`
	NoteTxID         string   `json:"note_txid"`
	NoteActionIndex  uint32   `json:"note_action_index"`
	NoteHeight       int64    `json:"note_height"`
	AmountZatoshis   uint64   `json:"amount_zatoshis"`
	NoteNullifier    string   `json:"note_nullifier,omitempty"`
	RecipientAddress string   `json:"recipient_address,omitempty"`
	Status           TxStatus `json:"status"`
}

type SpendConfirmedPayload struct {
	SpendEventPayload
	ConfirmedHeight       int64 `json:"confirmed_height"`
	RequiredConfirmations int64 `json:"required_confirmations"`
}

type SpendOrphanedPayload struct {
	SpendEventPayload
	OrphanedAtHeight int64 `json:"orphaned_at_height"`
}

type SpendUnconfirmedPayload struct {
	SpendEventPayload
	RollbackHeight          int64 `json:"rollback_height"`
	RequiredConfirmations   int64 `json:"required_confirmations,omitempty"`
	PreviousConfirmedHeight int64 `json:"previous_confirmed_height"`
}
