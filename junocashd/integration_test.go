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
