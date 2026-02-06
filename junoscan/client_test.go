package junoscan_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Abdullah1738/juno-sdk-go/junoscan"
	"github.com/Abdullah1738/juno-sdk-go/types"
)

func TestClient_UpsertWalletAndList(t *testing.T) {
	var (
		gotWalletID string
		gotUFVK     string
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/wallets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req map[string]string
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			gotWalletID = strings.TrimSpace(req["wallet_id"])
			gotUFVK = strings.TrimSpace(req["ufvk"])
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"wallets": []map[string]any{
					{"wallet_id": "hot", "created_at": time.Unix(1, 0).UTC()},
				},
			})
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	c, err := junoscan.New(srv.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.UpsertWallet(ctx, "hot", "ufvk123"); err != nil {
		t.Fatalf("UpsertWallet: %v", err)
	}
	if gotWalletID != "hot" {
		t.Fatalf("wallet_id=%q", gotWalletID)
	}
	if gotUFVK != "ufvk123" {
		t.Fatalf("ufvk=%q", gotUFVK)
	}

	wallets, err := c.ListWallets(ctx)
	if err != nil {
		t.Fatalf("ListWallets: %v", err)
	}
	if len(wallets) != 1 || wallets[0].WalletID != "hot" {
		t.Fatalf("unexpected wallets")
	}
}

func TestClient_ListWalletEvents(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/wallets/hot/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Query().Get("cursor") != "7" {
			http.Error(w, "bad cursor", http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("limit") != "123" {
			http.Error(w, "bad limit", http.StatusBadRequest)
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"events": []map[string]any{
				{
					"id":         8,
					"kind":       string(types.WalletEventKindDepositEvent),
					"height":     100,
					"payload":    json.RawMessage(`{"txid":"deadbeef"}`),
					"created_at": time.Unix(2, 0).UTC(),
				},
			},
			"next_cursor": 9,
		})
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	c, err := junoscan.New(srv.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page, err := c.ListWalletEvents(ctx, "hot", 7, 123)
	if err != nil {
		t.Fatalf("ListWalletEvents: %v", err)
	}
	if page.NextCursor != 9 || len(page.Events) != 1 {
		t.Fatalf("unexpected page")
	}
	if page.Events[0].Kind != types.WalletEventKindDepositEvent {
		t.Fatalf("kind=%q", page.Events[0].Kind)
	}
}

func TestClient_HTTPErrorIncludesStatusCode(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusBadRequest)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	c, err := junoscan.New(srv.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = c.Health(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	var he *junoscan.HTTPError
	if !errors.As(err, &he) {
		t.Fatalf("expected HTTPError, got %T", err)
	}
	if he.StatusCode != http.StatusBadRequest {
		t.Fatalf("status=%d", he.StatusCode)
	}
}

func TestClient_ListWalletNotesPage(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/wallets/hot/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		q := r.URL.Query()
		if q.Get("spent") != "false" {
			http.Error(w, "bad spent", http.StatusBadRequest)
			return
		}
		if q.Get("limit") != "200" {
			http.Error(w, "bad limit", http.StatusBadRequest)
			return
		}
		if q.Get("min_value_zat") != "1234" {
			http.Error(w, "bad min_value_zat", http.StatusBadRequest)
			return
		}
		if q.Get("cursor") != "11:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa:0" {
			http.Error(w, "bad cursor", http.StatusBadRequest)
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"notes": []map[string]any{
				{
					"txid":              "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
					"action_index":      0,
					"height":            12,
					"value_zat":         5000,
					"note_nullifier":    "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
					"recipient_address": "jtest1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqp4f3t7",
					"created_at":        time.Unix(3, 0).UTC(),
				},
			},
			"next_cursor": "12:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb:0",
		})
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	c, err := junoscan.New(srv.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page, err := c.ListWalletNotesPage(ctx, "hot", junoscan.ListWalletNotesOptions{
		OnlyUnspent: true,
		Limit:       200,
		MinValueZat: 1234,
		Cursor:      "11:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa:0",
	})
	if err != nil {
		t.Fatalf("ListWalletNotesPage: %v", err)
	}
	if len(page.Notes) != 1 {
		t.Fatalf("notes=%d", len(page.Notes))
	}
	if page.NextCursor == "" {
		t.Fatalf("expected next_cursor")
	}
}
