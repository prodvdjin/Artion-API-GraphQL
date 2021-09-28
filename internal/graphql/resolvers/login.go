package resolvers

import (
	"artion-api-graphql/internal/auth"
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"time"
)

// GenerateLoginChallenge generate login challenge, to be signed by client and passed into login mutation
func (rs *RootResolver) GenerateLoginChallenge() (string, error) {
	return auth.GetAuthChallenge(), nil
}

// GenerateBearer use signed login challenge to obtain Bearer token
func (rs *RootResolver) GenerateBearer(args struct{ user common.Address; challenge string; signature string }) (string, error) {
	timeInChallenge, err := time.Parse("2006-01-02-15:04:05.999999999", args.challenge[28:])
	if err != nil {
		return "", fmt.Errorf("unable to parse challenge; %s", err)
	}
	timeFromChallenge := time.Now().Sub(timeInChallenge).Minutes()
	if timeFromChallenge < 0 || timeFromChallenge > 60 {
		return "", fmt.Errorf("challenge is expired")
	}



}
