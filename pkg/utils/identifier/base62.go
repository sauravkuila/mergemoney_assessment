package identifier

import (
	"fmt"
	"time"

	"math/rand"

	"github.com/jxskiss/base62"
)

type base62Engine struct{}

func getBase62Obj() IdentifierItf {
	return &base62Engine{}
}

func (obj *base62Engine) GetUniqueId(appendData string) string {
	// Use current Unix timestamp in nanoseconds and a random value to reduce collision risk
	timestamp := time.Now().UnixNano()
	randomValue := rand.Int63()

	// Combine timestamp and random value into a single number
	combined := timestamp ^ randomValue

	// Encode the combined number to Base62
	return base62.EncodeToString([]byte(fmt.Sprintf("%d", combined)))
}

func GenerateUniqueBase62ID(appendData string) string {
	// Use current Unix timestamp in nanoseconds and a random value to reduce collision risk
	timestamp := time.Now().UnixNano()
	randomValue := rand.Int63()

	// Combine timestamp and random value into a single number
	combined := timestamp ^ randomValue

	// Encode the combined number to Base62
	return base62.EncodeToString([]byte(fmt.Sprintf("%d", combined)))
}
