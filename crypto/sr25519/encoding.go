package sr25519

import (
	"github.com/tendermint/tendermint/internal/jsontypes"
	tmjson "github.com/tendermint/tendermint/libs/json"
)

const (
	privKeyName = "tendermint/PrivKeySr25519"
	pubKeyName  = "tendermint/PubKeySr25519"
)

func init() {
	tmjson.RegisterType(PubKey{}, PubKeyName)
	tmjson.RegisterType(PrivKey{}, PrivKeyName)

	jsontypes.MustRegister(PubKey{})
	jsontypes.MustRegister(PrivKey{})
	crypto.RegisterPubKeyType(pubKeyName, PubKey(nil))
	crypto.RegisterPrivKeyType(privKeyName, PrivKey{})
}

//go:generate -command gen go run github.com/tendermint/tendermint/scripts/tmjson
//go:generate gen -output generated.go -pkg sr25519 PubKey=tendermint/PubKeySr25519
