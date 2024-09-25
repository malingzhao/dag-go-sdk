package tx

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"log"
)

// uint16ToLowBytes extracts the low byte of each uint16 element and converts to a byte slice.
func uint16ToLowBytes(data []uint16) []byte {
	var lowBytes []byte
	for _, v := range data {
		lowBytes = append(lowBytes, byte(v&0xFF)) // Only keep the lower byte
	}
	return lowBytes
}

// utf8Length encodes the length of a string as a variable-length integer
// with specific bits set to denote UTF-8 and additional bytes.
func utf8Length(value int) []byte {
	var buffer []uint16

	if value>>6 == 0 {
		// Use 1 byte
		buffer = append(buffer, uint16(value|0x80)) // Set bit 8.
	} else if value>>13 == 0 {
		// Use 2 bytes
		buffer = append(buffer, uint16(value|0x40|0x80)) // Set bit 7 and 8.
		buffer = append(buffer, uint16(value>>6))
	} else if value>>20 == 0 {
		// Use 3 bytes
		buffer = append(buffer, uint16(value|0x40|0x80)) // Set bit 7 and 8.
		buffer = append(buffer, uint16((value>>6)|0x80)) // Set bit 8.
		buffer = append(buffer, uint16(value>>13))
	} else if value>>27 == 0 {
		// Use 4 bytes
		buffer = append(buffer, uint16(value|0x40|0x80))  // Set bit 7 and 8.
		buffer = append(buffer, uint16((value>>6)|0x80))  // Set bit 8.
		buffer = append(buffer, uint16((value>>13)|0x80)) // Set bit 8.
		buffer = append(buffer, uint16(value>>20))
	} else {
		// Use 5 bytes
		buffer = append(buffer, uint16(value|0x40|0x80))  // Set bit 7 and 8.
		buffer = append(buffer, uint16((value>>6)|0x80))  // Set bit 8.
		buffer = append(buffer, uint16((value>>13)|0x80)) // Set bit 8.
		buffer = append(buffer, uint16((value>>20)|0x80)) // Set bit 8.
		buffer = append(buffer, uint16(value>>27))
	}

	return uint16ToLowBytes(buffer)
}

const MIN_SALT = 100000 // 设置最小盐值

// sha256Hash computes the SHA-256 hash of the input hex string.
func sha256Hash(serializedTx string) string {
	// Decode the hex string into bytes
	data, err := hex.DecodeString(serializedTx)
	if err != nil {
		fmt.Println("Error decoding hex:", err)
		return ""
	}

	// Compute the SHA-256 hash
	hash := sha256.Sum256(data)
	// Convert the hash to a hex string and return it
	return hex.EncodeToString(hash[:])
}

// kryoSerialize serializes the given message according to the specified rules.
func kryoSerialize(msg string, setReferences bool) string {
	s := utf8Length(len(msg) + 1)
	ext := hex.EncodeToString([]byte(s))

	prefix := "03"
	if setReferences {
		prefix += "01"
	}

	prefix = prefix + ext
	log.Println("the prefix is ", prefix)
	coded := hex.EncodeToString([]byte(msg))

	return prefix + coded
}

// sign signs the given message with the provided private key.
func sign(privateKeyHex string, msg string) (string, error) {
	// Calculate SHA-512 hash of the message
	hash := sha512.Sum512([]byte(msg))

	log.Println("the hash is ", hex.EncodeToString(hash[:]))
	// Decode the private key from hex
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("error decoding private key: %v", err)
	}

	// Import the private key
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	// Sign the hash
	signature := ecdsa.Sign(privKey, hash[:])

	return hex.EncodeToString(signature.Serialize()), nil
}
