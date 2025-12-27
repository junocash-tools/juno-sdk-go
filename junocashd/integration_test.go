//go:build integration

package junocashd_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Abdullah1738/juno-sdk-go/internal/testutil"
	"github.com/Abdullah1738/juno-sdk-go/junocashd"
)

func TestClient_GetBlockchainInfo_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	r, err := testutil.StartJunocashd(ctx, testutil.JunocashdConfig{})
	if err != nil {
		if errors.Is(err, testutil.ErrJunocashdNotFound) {
			t.Skip("junocashd not found in PATH")
		}
		t.Fatalf("StartJunocashd: %v", err)
	}
	defer func() { _ = r.Stop(context.Background()) }()

	cli := junocashd.New(r.RPCURL, r.RPCUser, r.RPCPassword)

	info, err := cli.GetBlockchainInfo(ctx)
	if err != nil {
		t.Fatalf("GetBlockchainInfo: %v", err)
	}
	if info.Chain != "regtest" {
		t.Fatalf("chain=%q want regtest", info.Chain)
	}
}

func TestClient_GetBlockVerbose_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	r, err := testutil.StartJunocashd(ctx, testutil.JunocashdConfig{})
	if err != nil {
		if errors.Is(err, testutil.ErrJunocashdNotFound) {
			t.Skip("junocashd not found in PATH")
		}
		t.Fatalf("StartJunocashd: %v", err)
	}
	defer func() { _ = r.Stop(context.Background()) }()

	cli := junocashd.New(r.RPCURL, r.RPCUser, r.RPCPassword)
	best, err := cli.GetBestBlockHash(ctx)
	if err != nil {
		t.Fatalf("GetBestBlockHash: %v", err)
	}

	block, err := cli.GetBlockVerbose(ctx, best)
	if err != nil {
		t.Fatalf("GetBlockVerbose: %v", err)
	}
	if block.Hash != best {
		t.Fatalf("hash=%q want %q", block.Hash, best)
	}
	if block.Height != 0 {
		t.Fatalf("height=%d want 0", block.Height)
	}
	if len(block.Tx) == 0 {
		t.Fatalf("expected at least one tx in genesis block")
	}
}
