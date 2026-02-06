package junoscan

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Option func(*Client)

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

func New(baseURL string, opts ...Option) (*Client, error) {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return nil, errors.New("junoscan: base url required")
	}
	u, err := url.Parse(baseURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("junoscan: invalid base url %q", baseURL)
	}

	c := &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	return c, nil
}

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	body := strings.TrimSpace(e.Body)
	if body == "" {
		return fmt.Sprintf("junoscan: http %d", e.StatusCode)
	}
	return fmt.Sprintf("junoscan: http %d: %s", e.StatusCode, body)
}

func (c *Client) Health(ctx context.Context) (HealthResponse, error) {
	var resp HealthResponse
	if err := c.doJSON(ctx, http.MethodGet, "/v1/health", nil, &resp); err != nil {
		return HealthResponse{}, err
	}
	return resp, nil
}

func (c *Client) UpsertWallet(ctx context.Context, walletID, ufvk string) error {
	walletID = strings.TrimSpace(walletID)
	ufvk = strings.TrimSpace(ufvk)
	if walletID == "" || ufvk == "" {
		return errors.New("junoscan: wallet_id and ufvk required")
	}

	var resp struct {
		Status string `json:"status"`
	}
	if err := c.doJSON(ctx, http.MethodPost, "/v1/wallets", walletRequest{WalletID: walletID, UFVK: ufvk}, &resp); err != nil {
		return err
	}
	if strings.ToLower(strings.TrimSpace(resp.Status)) != "ok" {
		return errors.New("junoscan: unexpected response")
	}
	return nil
}

func (c *Client) ListWallets(ctx context.Context) ([]Wallet, error) {
	var resp struct {
		Wallets []Wallet `json:"wallets"`
	}
	if err := c.doJSON(ctx, http.MethodGet, "/v1/wallets", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Wallets, nil
}

func (c *Client) ListWalletEvents(ctx context.Context, walletID string, cursor int64, limit int) (WalletEventsPage, error) {
	walletID = strings.TrimSpace(walletID)
	if walletID == "" {
		return WalletEventsPage{}, errors.New("junoscan: wallet_id required")
	}
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	path := fmt.Sprintf("/v1/wallets/%s/events?cursor=%d&limit=%d", url.PathEscape(walletID), cursor, limit)
	var resp WalletEventsPage
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return WalletEventsPage{}, err
	}
	return resp, nil
}

func (c *Client) ListWalletNotes(ctx context.Context, walletID string, onlyUnspent bool) ([]WalletNote, error) {
	page, err := c.ListWalletNotesPage(ctx, walletID, ListWalletNotesOptions{
		OnlyUnspent: onlyUnspent,
		Limit:       1000,
	})
	if err != nil {
		return nil, err
	}
	return page.Notes, nil
}

func (c *Client) ListWalletNotesPage(ctx context.Context, walletID string, opts ListWalletNotesOptions) (WalletNotesPage, error) {
	walletID = strings.TrimSpace(walletID)
	if walletID == "" {
		return WalletNotesPage{}, errors.New("junoscan: wallet_id required")
	}

	spentParam := "false"
	if !opts.OnlyUnspent {
		spentParam = "true"
	}

	limit := opts.Limit
	if limit <= 0 {
		limit = 1000
	}
	if limit > 1000 {
		limit = 1000
	}

	params := url.Values{}
	params.Set("spent", spentParam)
	params.Set("limit", fmt.Sprintf("%d", limit))
	if opts.MinValueZat > 0 {
		params.Set("min_value_zat", fmt.Sprintf("%d", opts.MinValueZat))
	}
	cursor := strings.TrimSpace(opts.Cursor)
	if cursor != "" {
		params.Set("cursor", cursor)
	}

	path := fmt.Sprintf("/v1/wallets/%s/notes?%s", url.PathEscape(walletID), params.Encode())
	var resp WalletNotesPage
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return WalletNotesPage{}, err
	}
	return resp, nil
}

func (c *Client) OrchardWitness(ctx context.Context, anchorHeight *int64, positions []uint32) (OrchardWitnessResponse, error) {
	if len(positions) == 0 {
		return OrchardWitnessResponse{}, errors.New("junoscan: positions required")
	}
	req := WitnessRequest{
		AnchorHeight: anchorHeight,
		Positions:    positions,
	}
	var resp OrchardWitnessResponse
	if err := c.doJSON(ctx, http.MethodPost, "/v1/orchard/witness", req, &resp); err != nil {
		return OrchardWitnessResponse{}, err
	}
	return resp, nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, in any, out any) error {
	if ctx == nil {
		ctx = context.Background()
	}

	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return errors.New("junoscan: marshal request")
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return errors.New("junoscan: build request")
	}
	req.Header.Set("Accept", "application/json")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	const maxBody = 1 << 20
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, maxBody))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &HTTPError{StatusCode: resp.StatusCode, Body: string(raw)}
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return errors.New("junoscan: invalid json response")
	}
	return nil
}
