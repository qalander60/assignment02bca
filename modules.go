package assignment02bca

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// WriteLastHash writes the hash string to a file named "last_hash.txt"
func WriteLastHash(hash string) error {
	return os.WriteFile("last_hash.txt", []byte(hash), 0644)
}

// ReadLastHash reads the hash string from the file named "last_hash.txt"
func ReadLastHash() (string, error) {
	data, err := os.ReadFile("last_hash.txt")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// structure for the transaction
type Transaction struct {
	TransactionID              string
	SenderBlockchainAddress    string
	RecipientBlockchainAddress string
	Value                      float32
}

// updated block structure with timestamp
type Block struct {
	transactions []*Transaction
	nonce        int
	previousHash string
	currentHash  string
	id           int
	timestamp    time.Time
	next         *Block
}

type Chain struct {
	head            *Block
	transactionPool []*Transaction
	length          int
}

// Hash function that serializes the block, calculates the hash, and returns it as a string
func CalculateHash(stringToHash string) string {
	hash := sha256.Sum256([]byte(stringToHash))
	return fmt.Sprintf("%x", hash[:])
}

// function for the new transaction
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	transactionData := fmt.Sprintf("%s%s%f", sender, recipient, value)
	return &Transaction{
		TransactionID:              CalculateHash(transactionData),
		SenderBlockchainAddress:    sender,
		RecipientBlockchainAddress: recipient,
		Value:                      value,
	}
}

// Method to add a transaction to the central transaction pool
func (bc *Chain) AddTransactionToPool(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

// ProofOfWork adjusts the nonce until the hash meets the difficulty criteria
func (b *Block) ProofOfWork(difficulty int) {
	prefix := strings.Repeat("a", difficulty)
	for {
		blockData := fmt.Sprintf("%v%d%s%d", b.transactions, b.nonce, b.previousHash, b.id)
		b.currentHash = CalculateHash(blockData)
		if strings.HasPrefix(b.currentHash, prefix) {
			break
		}
		b.nonce++
	}
}

// function to create a new block with transactions from the pool
func (bc *Chain) NewBlock() *Block {
	if len(bc.transactionPool) == 0 {
		fmt.Println("No transactions available in the pool.")
		return nil
	}

	id := bc.length

	// Pop up to 2 transactions from the transaction pool
	transactions := []*Transaction{}
	if len(bc.transactionPool) > 0 {
		transactions = append(transactions, bc.transactionPool[0])
		bc.transactionPool = bc.transactionPool[1:]
	}
	if len(bc.transactionPool) > 0 {
		transactions = append(transactions, bc.transactionPool[0])
		bc.transactionPool = bc.transactionPool[1:]
	}

	// Create new block with timestamp
	newBlock := &Block{
		transactions: transactions,
		nonce:        0,
		previousHash: "",
		currentHash:  "",
		id:           id,
		timestamp:    time.Now(),
		next:         nil,
	}

	// If this is the first block, make it the head (genesis block)
	if bc.head == nil {
		bc.head = newBlock
		newBlock.ProofOfWork(2)
		fmt.Println("Genesis Block created.")
		bc.length++
		err := WriteLastHash(newBlock.currentHash)
		if err != nil {
			fmt.Println("Error writing last hash to file:", err)
		}
		return newBlock
	}

	// Find the last block in the chain
	currBlock := bc.head
	for currBlock.next != nil {
		currBlock = currBlock.next
	}

	// Link the new block and calculate proof of work
	newBlock.previousHash = currBlock.currentHash
	newBlock.ProofOfWork(2)
	currBlock.next = newBlock
	bc.length++
	err := WriteLastHash(newBlock.currentHash)
	if err != nil {
		fmt.Println("Error writing last hash to file:", err)
	}
	return newBlock
}

// function to list all blocks
func (bc *Chain) ListBlocks() {
	if bc.head == nil {
		fmt.Println("Blockchain is empty.")
		return
	}
	currBlock := bc.head
	for currBlock != nil {
		fmt.Println("--------------------------------------------------")
		fmt.Printf("Block ID       : %d\n", currBlock.id)
		fmt.Printf("Nonce          : %d\n", currBlock.nonce)
		fmt.Printf("Previous Hash  : %s\n", currBlock.previousHash)
		fmt.Printf("Current Hash   : %s\n", currBlock.currentHash)

		// Convert transactions to JSON format
		jsonTransactions, err := json.MarshalIndent(currBlock.transactions, "", "  ")
		if err != nil {
			fmt.Println("Error converting transactions to JSON:", err)
		} else {
			fmt.Printf("Transactions   : %s\n", string(jsonTransactions))
		}

		fmt.Println("--------------------------------------------------")
		fmt.Println()
		currBlock = currBlock.next
	}
}
