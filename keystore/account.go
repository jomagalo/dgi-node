package keystore

import "github.com/jomagalo/dgi-node/common"

type Account struct {
	Address common.Address `json:"address"` // Ethereum account address derived from the key
	URL     URL            `json:"url"`     // Optional resource locator within a backend
}
