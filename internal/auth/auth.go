package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"time"
)

func GetAuthChallenge() string {
	tm := time.Now().Format("2006-01-02-15:04:05.999999999")
	challenge := fmt.Sprintf("I am signing into Artion at %s", tm)
	return challenge
}

type Bearer string {
	challenge string
	signature
}

func VerifyBearer(bearer string) {
	bearerBytes := hexutil.Decode(bearer)
	signature := bearerBytes[len(bearerBytes)-65:]

	json.Unmarshal(bearerBytes)




}


// signature has always 65 bytes


func checkChallengeSignature(user common.Address, challenge string, signatureHex string) bool {
	hash := crypto.Keccak256Hash([]byte(challenge))
	signature, err := hexutil.Decode(signatureHex)

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}
	signatureAddress := crypto.PubkeyToAddress(*sigPublicKeyECDSA)

	return bytes.Equal(signatureAddress.Bytes(), user.Bytes())
}
