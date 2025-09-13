package identifier

type IdentifierItf interface {
	//returns a unique id everytime with appended argument
	GetUniqueId(string) string
}

func GetIdentifierItf(engine IdentifierType) IdentifierItf {
	if engine == IDENTIFIER_BASE62 {
		return getBase62Obj()
	}
	return getSnowflakeObj()
}
