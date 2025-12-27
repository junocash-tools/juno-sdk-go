# juno-sdk-go

Go SDK and shared types for the Juno Cash toolchain.

Includes shared types (events, notes, TxPlan, error codes) for withdrawals, sweeps, and wallet tiering.

Status: work in progress.

## Development

Prereqs:

- Go 1.22+
- `junocashd` and `junocash-cli` in `PATH` (for integration/e2e)

Commands:

- Unit tests: `make test`
- Integration tests (starts a local `junocashd -regtest`): `make test-integration`
- E2E tests (daemon + CLI): `make test-e2e`

Set `JUNO_TEST_LOG=1` to show `junocashd` logs while running tests (it also enables `go test -v` via the Makefile).

## Packages

- `junocashd`: typed JSON-RPC client helpers for `junocashd` (blocks, headers, tx broadcast)
- `types`: shared payload types (TxPlan, DepositEvent, ChainCursor, stable error codes)
