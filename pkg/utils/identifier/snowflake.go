package identifier

import (
	"fmt"

	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake = nil

type snowflakeEngine struct {
	sf *sonyflake.Sonyflake
}

func getSnowflakeObj() IdentifierItf {
	return &snowflakeEngine{
		sf: sf,
	}
}

func (obj *snowflakeEngine) GetUniqueId(appendData string) string {
	// Create a new sonyflake instance with default settings. done only on first call
	if obj.sf == nil {
		obj.sf = sonyflake.NewSonyflake(sonyflake.Settings{})
	}

	// Generate a unique ID
	id, err := obj.sf.NextID()
	if err != nil {
		return ""
	}

	// Return the ID as a string (you can use base62 encoding for shortness)
	return fmt.Sprintf("%s_%d", appendData, id)
}

func GenerateUniqueSnowflakeID(appendData string) string {
	// Create a new sonyflake instance with default settings. done only on first call
	if sf == nil {
		sf = sonyflake.NewSonyflake(sonyflake.Settings{})
	}

	// Generate a unique ID
	id, err := sf.NextID()
	if err != nil {
		return ""
	}

	// Return the ID as a string (you can use base62 encoding for shortness)
	return fmt.Sprintf("%s_%d", appendData, id)
}
