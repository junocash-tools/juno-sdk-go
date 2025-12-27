package testutil

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/Abdullah1738/juno-sdk-go/junocashd"
)

var ErrJunocashdNotFound = errors.New("junocashd not found")

type JunocashdConfig struct {
	JunocashdPath string
	RPCUser       string
	RPCPassword   string
}

type RunningJunocashd struct {
	cmd         *exec.Cmd
	Datadir     string
	RPCURL      string
	RPCPort     int
	RPCUser     string
	RPCPassword string

	stopOnce sync.Once
}

func StartJunocashd(ctx context.Context, cfg JunocashdConfig) (*RunningJunocashd, error) {
	bin := cfg.JunocashdPath
	if bin == "" {
		p, err := exec.LookPath("junocashd")
		if err != nil {
			return nil, ErrJunocashdNotFound
		}
		bin = p
	}

	rpcUser := cfg.RPCUser
	if rpcUser == "" {
		rpcUser = "rpcuser"
	}
	rpcPass := cfg.RPCPassword
	if rpcPass == "" {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return nil, fmt.Errorf("rand: %w", err)
		}
		rpcPass = hex.EncodeToString(b)
	}

	rpcPort, err := freePort()
	if err != nil {
		return nil, err
	}
	p2pPort, err := freePort()
	if err != nil {
		return nil, err
	}

	datadir, err := os.MkdirTemp("", "junocashd-regtest-*")
	if err != nil {
		return nil, fmt.Errorf("mkdtemp: %w", err)
	}

	args := []string{
		"-regtest",
		"-server=1",
		"-daemon=0",
		"-listen=0",
		"-printtoconsole=1",
		fmt.Sprintf("-datadir=%s", datadir),
		fmt.Sprintf("-rpcbind=%s", "127.0.0.1"),
		fmt.Sprintf("-rpcallowip=%s", "127.0.0.1"),
		fmt.Sprintf("-rpcport=%d", rpcPort),
		fmt.Sprintf("-rpcuser=%s", rpcUser),
		fmt.Sprintf("-rpcpassword=%s", rpcPass),
		fmt.Sprintf("-port=%d", p2pPort),
	}

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = datadir
	logs := io.Discard
	if os.Getenv("JUNO_TEST_LOG") != "" {
		logs = os.Stdout
	}
	cmd.Stdout = logs
	cmd.Stderr = logs

	if err := cmd.Start(); err != nil {
		_ = os.RemoveAll(datadir)
		return nil, fmt.Errorf("start junocashd: %w", err)
	}

	r := &RunningJunocashd{
		cmd:         cmd,
		Datadir:     datadir,
		RPCURL:      fmt.Sprintf("http://127.0.0.1:%d", rpcPort),
		RPCPort:     rpcPort,
		RPCUser:     rpcUser,
		RPCPassword: rpcPass,
	}

	if err := r.waitForRPC(ctx, 25*time.Second); err != nil {
		_ = r.Stop(context.Background())
		return nil, err
	}
	return r, nil
}

func (r *RunningJunocashd) Stop(ctx context.Context) error {
	var err error
	r.stopOnce.Do(func() {
		if r.cmd == nil || r.cmd.Process == nil {
			err = nil
			return
		}

		// Prefer a graceful stop.
		if runtime.GOOS == "windows" {
			_ = r.cmd.Process.Kill()
		} else {
			_ = r.cmd.Process.Signal(syscall.SIGTERM)
		}

		done := make(chan struct{})
		go func() {
			_ = r.cmd.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-ctx.Done():
			_ = r.cmd.Process.Kill()
			_ = r.cmd.Wait()
		}

		_ = os.RemoveAll(r.Datadir)
	})
	return err
}

func (r *RunningJunocashd) waitForRPC(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cli := junocashd.New(r.RPCURL, r.RPCUser, r.RPCPassword)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		_, err := cli.GetBlockchainInfo(ctx)
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("junocashd rpc not ready: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

func freePort() (int, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("listen: %w", err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port, nil
}
