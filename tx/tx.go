package tx

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type Proof struct {
	Signature string `json:"signature"`
	Id        string `json:"id"`
}

type PostTxValue struct {
	Source      string           `json:"source"`
	Destination string           `json:"destination"`
	Amount      *big.Int         `json:"amount"`
	Fee         int64            `json:"fee"`
	Parent      AddressLastRefV2 `json:"parent"`
	Salt        *big.Int         `json:"salt"`
}

type Transaction struct {
	Value  PostTxValue `json:"value"`
	Proofs []Proof     `json:"proofs"`
}

type AddressLastRefV2 struct {
	Hash    string `json:"hash"`
	Ordinal int    `json:"ordinal"`
}

func (tx *Transaction) GetEncoded() string {
	var builder strings.Builder
	// Always 2 parents
	parentCount := "2"
	builder.WriteString(parentCount)

	// Source address and its length
	sourceAddress := tx.Value.Source
	builder.WriteString(strconv.Itoa(len(sourceAddress)))
	builder.WriteString(sourceAddress)

	// Destination address and its length
	destAddress := tx.Value.Destination
	builder.WriteString(strconv.Itoa(len(destAddress)))
	builder.WriteString(destAddress)

	// Amount as hex and its length
	amountHex := tx.Value.Amount.Text(16)
	builder.WriteString(strconv.Itoa(len(amountHex)))
	builder.WriteString(amountHex)

	// Parent hash and its length
	parentHash := tx.Value.Parent.Hash
	builder.WriteString(strconv.Itoa(len(parentHash)))
	builder.WriteString(parentHash)

	// Ordinal as string and its length
	ordinal := strconv.Itoa(tx.Value.Parent.Ordinal)
	builder.WriteString(strconv.Itoa(len(ordinal)))
	builder.WriteString(ordinal)

	// Fee as string and its length
	fee := fmt.Sprintf("%d", tx.Value.Fee)
	builder.WriteString(strconv.Itoa(len(fee)))
	builder.WriteString(fee)

	// Salt as hex and its length
	saltHex := tx.Value.Salt.Text(16)
	builder.WriteString(strconv.Itoa(len(saltHex)))
	builder.WriteString(saltHex)
	// Return the final encoded string
	return builder.String()
}

func (t *Transaction) Sign(privateKey string, unCompressKey string) error {
	encoded := t.GetEncoded()
	serializeTx := kryoSerialize(encoded, false)
	hash := sha256Hash(serializeTx)
	s, err := sign(privateKey, hash)
	if err != nil {
		return err
	}
	if unCompressKey[0:2] == "0x" {
		unCompressKey = unCompressKey[2:]
	}
	unCompressKey = unCompressKey[2:]
	t.Proofs = make([]Proof, 0)
	t.Proofs = append(t.Proofs, Proof{
		Signature: s,
		Id:        unCompressKey,
	})
	return nil
}

func NewTx(amount *big.Int, from string, to string, lastRef AddressLastRefV2, fee int64) (*Transaction, error) {
	tx := new(Transaction)
	tx.Value.Source = from
	tx.Value.Destination = to
	tx.Value.Amount = amount
	tx.Value.Fee = fee
	tx.Value.Parent = lastRef
	// 生成盐值
	var salt *big.Int

	tx.Value.Salt = salt
	if salt == nil {
		randomBytes := make([]byte, 6)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return nil, err
		}
		randomValue := new(big.Int).SetBytes(randomBytes)
		salt = new(big.Int).SetInt64(MIN_SALT)
		salt = salt.Add(salt, randomValue)
		tx.Value.Salt = salt
	}
	return tx, nil
}
