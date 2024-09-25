package tx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"log"
	"math/big"
	"testing"
)

func Test_Tx(t *testing.T) {
	tx, err := NewTx(new(big.Int).SetInt64(20000000), "DAG85Z8DiU44dbtUwsPW9crHadUT5CAryQbYVnFW", "DAG6n4UK5gncK6B2pYGq3NSwnxqi4EdcZJdadBhA", AddressLastRefV2{
		Hash:    "e116cdd6535915e1e0638151520272cc2c73c5f247f7429d06c4a3ca948c734e",
		Ordinal: 3,
	}, 10000000)

	if err != nil {
		log.Println("new Tx err ", err)
		return
	}
	err = tx.Sign("", "04025fa844d48aba12a6e9c53057e52c23c0fe336724f57ba0edc960e85a8b88d15a2ad8d7754fa72c7f6369eebdd17cb61ae604a8c7e4242a713680a3ebb005ef")
	marshal, err := json.Marshal(tx)
	if err != nil {
		log.Println("marshal err ", err)
		return
	}
	log.Println(string(marshal))

}

func Test_Priv(t *testing.T) {
	decodeString, _ := hex.DecodeString("96b3c2975ce2b436208a7a12428ed4139a30c310006e3b0791a9aba114872a9d")
	_, publicKey := btcec.PrivKeyFromBytes(decodeString)
	fmt.Println("the hex is ", hex.EncodeToString(publicKey.SerializeUncompressed()))
	//ecdsa.SignCompact()
}
