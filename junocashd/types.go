package junocashd

type BlockchainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int64   `json:"blocks"`
	Headers              int64   `json:"headers,omitempty"`
	BestBlockHash        string  `json:"bestblockhash,omitempty"`
	InitialBlockDownload bool    `json:"initialblockdownload,omitempty"`
	VerificationProgress float64 `json:"verificationprogress,omitempty"`
	Difficulty           float64 `json:"difficulty,omitempty"`
	ChainWork            string  `json:"chainwork,omitempty"`
	Pruned               bool    `json:"pruned,omitempty"`
	PruneHeight          int64   `json:"pruneheight,omitempty"`
	SizeOnDisk           int64   `json:"size_on_disk,omitempty"`
}

type BlockHeader struct {
	Hash              string `json:"hash"`
	Confirmations     int64  `json:"confirmations,omitempty"`
	Height            int64  `json:"height"`
	Time              int64  `json:"time"`
	PreviousBlockHash string `json:"previousblockhash,omitempty"`
	NextBlockHash     string `json:"nextblockhash,omitempty"`
}

type BlockVerbose struct {
	Hash              string   `json:"hash"`
	Confirmations     int64    `json:"confirmations,omitempty"`
	Height            int64    `json:"height"`
	Time              int64    `json:"time"`
	PreviousBlockHash string   `json:"previousblockhash,omitempty"`
	NextBlockHash     string   `json:"nextblockhash,omitempty"`
	Tx                []string `json:"tx"`
}
