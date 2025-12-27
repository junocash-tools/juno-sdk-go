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
