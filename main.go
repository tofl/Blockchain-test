package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"github.com/joho/godotenv"
)

type Block struct {
	Index int `json:"index"`
	Timestamp string `json:"timestamp"`
	BPM int `json:"bpm"`
	PrevHash string `json:"prevHash"`
	Hash string `json:"hash"`
}

var Blockchain []Block

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

func newBlock(previousBlock Block, BPM int) (Block, error) {
	var newBlock Block

	newBlock.Index     = previousBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.BPM       = BPM
	newBlock.PrevHash  = previousBlock.Hash
	newBlock.Hash      = calculateHash(newBlock)

	return newBlock, nil
}

func blockIsValid(block, previousBlock Block) bool {
	if previousBlock.Index+1 != block.Index {
		return false
	}

	if block.PrevHash != previousBlock.Hash {
		return false
	}

	if calculateHash(block) != block.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.\n", err)
	}

	go func() {
		startBlock := Block{
			Index    : 0,
			Timestamp: time.Now().String(),
			BPM      : 0,
			Hash     : "",
			PrevHash : "",
		}

		Blockchain = append(Blockchain, startBlock)
	}()

	log.Fatal(runServer())
}