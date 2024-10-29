package domain

type MetadataInfo struct {
	Length      int    `json:"length"`
	Name        string `json:"name"`
	PieceLength int    `json:"piece_length"`
	Pieces      []byte `json:"pieces"`
}

type Metadata struct {
	Announce  string       `json:"announce"`
	CreatedBy string       `json:"created_by"`
	Info      MetadataInfo `json:"info"`
}
