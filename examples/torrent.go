package main

// This program will read it's first argument,
// parse it as a bencoded torrent file and
// unmarshal it to stdout.

import (
	"fmt"
	"github.com/tumdum/bencoding"
	"os"
)

type Info struct {
	Length      int64  `bencoding:"length"`
	Name        string `bencoding:"name"`
	PieceLength int64  `bencoding:"piece length"`
}

type Torrent struct {
	Announce     string   `bencoding:"announce"`
	AnnounceList []interface{} `bencoding:"announce-list"`
	Comment      string   `bencoding:"comment"`
	CreationDate int64    `bencoding:"creation date"`
	Info         Info     `bencoding:"info"`
}

func usage() {
	fmt.Printf("usage: $%s file.torrent\n", os.Args[0])
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(0)
	}
	f, e := os.Open(os.Args[1])
	if e != nil {
		panic(e)
	}
	defer f.Close()
	d := bencoding.NewTorrentDecoder(f)
	var t Torrent
	if h, e := d.Decode(&t); e != nil {
    panic(e)
	} else {
		fmt.Printf("  torrent: %+v\n     hash: %v\n", t, h)
	}
}
