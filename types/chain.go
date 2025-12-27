package types

type ChainCursor struct {
	Height int64  `json:"height"`
	Hash   string `json:"hash"`
}

type ReorgEvent struct {
	From ChainCursor `json:"from"`
	To   ChainCursor `json:"to"`
}
