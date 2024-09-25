package address

import (
	"fmt"
	"testing"
)

func Test_GetAddressByPublicKey(t *testing.T) {
	addr := GetDagAddressFromPublicKey("04025fa844d48aba12a6e9c53057e52c23c0fe336724f57ba0edc960e85a8b88d15a2ad8d7754fa72c7f6369eebdd17cb61ae604a8c7e4242a713680a3ebb005ef")

	fmt.Println("addr is ", addr)
}

func Test_ValidateAddress(t *testing.T) {
	fmt.Println(ValidateDagAddress("DAG85Z8DzU44dbtUwsPW9crHadUT5CAryQbYVnFW"))
	fmt.Println(ValidateDagAddress("DAG85Z8DzU44dbtUwsPW9crHadUT5CAr0QbYVnFW"))
}
