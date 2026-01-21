package junoscan

import (
	"encoding/json"
	"time"

	"github.com/Abdullah1738/juno-sdk-go/types"
)

type HealthResponse struct {
	Status        string  `json:"status"`
	ScannedHeight *int64  `json:"scanned_height,omitempty"`
	ScannedHash   *string `json:"scanned_hash,omitempty"`
}

type Wallet struct {
	WalletID   string     `json:"wallet_id"`
	CreatedAt  time.Time  `json:"created_at"`
	DisabledAt *time.Time `json:"disabled_at,omitempty"`
}

type WalletEvent struct {
	ID        int64                 `json:"id"`
	Kind      types.WalletEventKind `json:"kind"`
	Height    int64                 `json:"height"`
	Payload   json.RawMessage       `json:"payload"`
	CreatedAt time.Time             `json:"created_at"`
}

type WalletEventsPage struct {
	Events     []WalletEvent `json:"events"`
	NextCursor int64         `json:"next_cursor"`
}

type WalletNote struct {
	TxID             string     `json:"txid"`
	ActionIndex      int32      `json:"action_index"`
	Height           int64      `json:"height"`
	Position         *int64     `json:"position,omitempty"`
	RecipientAddress string     `json:"recipient_address"`
	ValueZat         int64      `json:"value_zat"`
	NoteNullifier    string     `json:"note_nullifier"`
	PendingSpentTxID *string    `json:"pending_spent_txid,omitempty"`
	PendingSpentAt   *time.Time `json:"pending_spent_at,omitempty"`
	SpentHeight      *int64     `json:"spent_height,omitempty"`
	SpentTxID        *string    `json:"spent_txid,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type WitnessRequest struct {
	AnchorHeight *int64   `json:"anchor_height,omitempty"`
	Positions    []uint32 `json:"positions"`
}

type OrchardWitnessPath struct {
	Position uint32   `json:"position"`
	AuthPath []string `json:"auth_path"`
}

type OrchardWitnessResponse struct {
	Status       string               `json:"status"`
	AnchorHeight int64                `json:"anchor_height"`
	Root         string               `json:"root"`
	Paths        []OrchardWitnessPath `json:"paths"`
}

type walletRequest struct {
	WalletID string `json:"wallet_id"`
	UFVK     string `json:"ufvk"`
}
