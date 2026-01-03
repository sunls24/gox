package header

import (
	"fmt"

	"github.com/sunls24/gox/types"
)

var ContentTypeJson = types.NewPair("Content-Type", "application/json")

func Authorization(token string) types.Pair[string] {
	return types.NewPair("Authorization", fmt.Sprintf("Bearer %s", token))
}
