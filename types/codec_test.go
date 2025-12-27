package types_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Abdullah1738/juno-sdk-go/types"
)

func TestTxPlan_JSONRoundTrip(t *testing.T) {
	in := types.TxPlan{
		Version:  types.V1,
		Kind:     types.TxPlanKindSweep,
		WalletID: "hot",
		Inputs: []types.NoteRef{
			{TxID: "deadbeef", ActionIndex: 7},
		},
		Outputs: []types.TxOutput{
			{ToAddress: "j1test", AmountZatoshis: 12345, MemoHex: "00"},
		},
		FeeZatoshis: 1000,
		Metadata: map[string]any{
			"reason": "test",
		},
	}

	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out types.TxPlan
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("round-trip mismatch:\n  in=%#v\n out=%#v", in, out)
	}
}

func TestDepositEvent_JSONRoundTrip(t *testing.T) {
	in := types.DepositEvent{
		Version:          types.V1,
		WalletID:         "hot",
		AccountID:        "acct_123",
		DiversifierIndex: 1,
		TxID:             "deadbeef",
		Height:           100,
		ActionIndex:      0,
		AmountZatoshis:   5000,
		MemoHex:          "",
		Status: types.TxStatus{
			State:         types.TxStateConfirmed,
			Height:        100,
			Confirmations: 3,
		},
	}

	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out types.DepositEvent
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("round-trip mismatch:\n  in=%#v\n out=%#v", in, out)
	}
}

func TestChainCursor_JSONRoundTrip(t *testing.T) {
	in := types.ChainCursor{
		Height: 123,
		Hash:   "deadbeef",
	}

	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out types.ChainCursor
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("round-trip mismatch:\n  in=%#v\n out=%#v", in, out)
	}
}
