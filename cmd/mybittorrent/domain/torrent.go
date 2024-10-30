package domain

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/decoder"
	"io"
	"os"
)

type TorrentFile struct {
	FilePath string

	Announce  string      `json:"announce"`
	CreatedBy string      `json:"created_by"`
	Info      TorrentInfo `json:"info"`
}

type TorrentInfo struct {
	Length      int    `json:"length"`
	Name        string `json:"name"`
	PieceLength int    `json:"piece_length"`
	Pieces      []byte `json:"pieces"`
}

func NewTorrent(filePath string) (*TorrentFile, error) {
	decoder := decoder.NewBencodeDecoder()

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	decoded, _, err := decoder.DecodeBencode([]rune(string(data)))
	if err != nil {
		return nil, err
	}

	m := decoded.(map[string]interface{})
	info := m["info"].(map[string]interface{})

	return &TorrentFile{
		FilePath:  filePath,
		Announce:  m["announce"].(string),
		CreatedBy: m["created by"].(string),
		Info: TorrentInfo{
			Length:      info["length"].(int),
			Name:        info["name"].(string),
			PieceLength: info["piece length"].(int),
			Pieces:      []byte(info["pieces"].(string)),
		},
	}, nil
}

func (torrent *TorrentFile) InfoHash() ([20]byte, error) {
	data, err := os.ReadFile(torrent.FilePath)
	if err != nil {
		return [20]byte{}, err
	}

	// 4:infod<INFO_CONTENTS>e
	infoStart := bytes.Index(data, []byte("4:info")) + 6
	if infoStart < 0 {
		return [20]byte{}, fmt.Errorf("TorrentFile.info: no info in torrent file")
	}
	return sha1.Sum(data[infoStart : len(data)-1]), nil
}
