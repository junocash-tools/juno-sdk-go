//go:build e2e && docker

package junocashd_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Abdullah1738/juno-sdk-go/internal/testutil/containers"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	jd, err := containers.StartJunocashd(ctx)
	if err != nil {
		log.Printf("start junocashd container: %v", err)
		os.Exit(1)
	}

	os.Setenv("JUNO_TEST_DOCKER", "1")
	os.Setenv("JUNO_TEST_RPC_URL", jd.RPCURL)
	os.Setenv("JUNO_TEST_RPC_USER", jd.RPCUser)
	os.Setenv("JUNO_TEST_RPC_PASS", jd.RPCPassword)
	os.Setenv("JUNO_TEST_JUNOCASHD_CONTAINER", jd.ContainerID)
	os.Setenv("JUNO_TEST_JUNOCASHD_DATADIR", "/data")
	os.Setenv("JUNO_TEST_JUNOCASHD_RPC_PORT", "8232")

	code := m.Run()

	termCtx, termCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer termCancel()
	if err := jd.Terminate(termCtx); err != nil {
		log.Printf("terminate container: %v", err)
	}

	os.Exit(code)
}
