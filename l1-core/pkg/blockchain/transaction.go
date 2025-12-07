package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type Transaction struct {
	ID        string                 `json:"id"`
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	Amount    float64                `json:"amount"`
	Fee       float64                `json:"fee"`
	Timestamp int64                  `json:"timestamp"`
	Signature string                 `json:"signature"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

func NewTransaction(from, to string, amount, fee float64, data map[string]interface{}) *Transaction {
	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Fee:       fee,
		Timestamp: time. Now().Unix(),
		Data:      data,
	}
	tx.ID = tx.CalculateHash()
	return tx
}

func (tx *Transaction) CalculateHash() string {
	data, _ := json.Marshal(tx)
	hash := sha256.Sum256(data)
	return fmt. Sprintf("%x", hash)
}

func (tx *Transaction) IsValid() bool {
	return tx.ID != "" && tx.From != "" && tx.To != "" && tx.Amount > 0
}
