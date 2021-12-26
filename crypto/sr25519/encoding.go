package sr25519

import (
	"github.com/tendermint/tendermint/internal/jsontypes"
	tmjson "github.com/tendermint/tendermint/libs/json"
)

// Generate JSON encoding wrappers for the types in this package.
//go:generate -command gen go run github.com/creachadair/misctools/tagtype@latest
//go:generate gen -output generated.go

func (PubKey) jsonWrapperTag() string  { return PubKeyName }
func (PrivKey) jsonWrapperTag() string { return PrivKeyName }

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
