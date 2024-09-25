package address

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/mr-tron/base58"
	"strings"
)

const PKCSPrefix = "3056301006072a8648ce3d020106052b8104000a034200"

func GetDagAddressFromPublicKey(publicKeyHex string) string {
	// 检查公钥是否缺少 '04' 前缀
	if len(publicKeyHex) == 128 {
		publicKeyHex = "04" + publicKeyHex
	}

	publicKeyHex = PKCSPrefix + publicKeyHex

	publicKeyBytes, _ := hex.DecodeString(publicKeyHex)
	sha256Hash := sha256.Sum256(publicKeyBytes)

	hash := base58.Encode(sha256Hash[:])

	end := hash[len(hash)-36:]

	sum := 0
	for _, char := range end {
		if char >= '0' && char <= '9' {
			sum += int(char - '0')
		}
	}
	par := sum % 9

	return fmt.Sprintf("DAG%d%s", par, end)
}

// validateDagAddress checks if the given address is a valid DAG address.
func ValidateDagAddress(address string) bool {
	if address == "" {
		return false
	}

	validLen := len(address) == 40
	validPrefix := strings.HasPrefix(address, "DAG")

	if len(address) < 4 {
		return false
	}
	par := address[3] - '0' // Convert char to int
	validParity := par >= 0 && par < 10

	// Using regex to match base58 characters
	_, err := base58.Decode(address[4:])
	if err != nil {
		return false
	}
	return validLen && validPrefix && validParity
}
