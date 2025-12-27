//go:build e2e

package junocashd_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/Abdullah1738/juno-sdk-go/internal/testutil"
	"github.com/Abdullah1738/juno-sdk-go/junocashd"
)

func TestJunocashdAndCLI_E2E(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := exec.LookPath("junocash-cli")
	if err != nil {
		t.Skip("junocash-cli not found in PATH")
	}

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

	out, err := exec.CommandContext(
		ctx,
		"junocash-cli",
		"-regtest",
		"-datadir="+r.Datadir,
		"-rpcuser="+r.RPCUser,
		"-rpcpassword="+r.RPCPassword,
		"-rpcport="+fmt.Sprint(r.RPCPort),
		"getblockchaininfo",
	).CombinedOutput()
	if err != nil {
		t.Fatalf("junocash-cli getblockchaininfo: %v\n%s", err, string(out))
	}

	var cliInfo struct {
		Chain string `json:"chain"`
	}
	if err := json.Unmarshal(out, &cliInfo); err != nil {
		t.Fatalf("unmarshal junocash-cli output: %v\n%s", err, string(out))
	}
	if cliInfo.Chain != "regtest" {
		t.Fatalf("cli chain=%q want regtest", cliInfo.Chain)
	}
}
