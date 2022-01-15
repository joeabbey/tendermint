package sr25519

import (
	"github.com/tendermint/tendermint/internal/jsontypes"
	tmjson "github.com/tendermint/tendermint/libs/json"
)

//go:generate -command gen go run github.com/tendermint/tendermint/scripts/tmjson
//go:generate gen -output generated.go -pkg sr25519 -m -prefix tendermint/ PubKey=+PubKeySr25519 PrivKey=+PrivKeySr25519

const (
	PrivKeyName = "tendermint/PrivKeySr25519"
	PubKeyName  = "tendermint/PubKeySr25519"
)

func init() {
	tmjson.RegisterType(PubKey{}, PubKeyName)
	tmjson.RegisterType(PrivKey{}, PrivKeyName)

	jsontypes.MustRegister(PubKey{})
	jsontypes.MustRegister(PrivKey{})
	crypto.RegisterPubKeyType(pubKeyName, PubKey(nil))
	crypto.RegisterPrivKeyType(privKeyName, PrivKey{})
}
