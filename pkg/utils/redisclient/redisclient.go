package redisclient

import "context"

type RedisClientItf interface {
	//set a value into redis with an expiry
	//	depending on params similar redis command
	//		writeIfNotSet: true -> SET key value EX timeout NX
	//	params
	//		key -> string
	//		value -> interface
	//		timeout -> int | value in millisecond
	//		writeIfNotSet -> bool
	//	interfaces can be directly sent as reference arguments
	//		eg. usage
	//		type Person struct{
	//			Name string,
	//			Age int,
	//		}
	//		var p1 Person{ Name: "John Doe", Age: 25 }
	//		ctx := context.Background()
	//		redisObj := utils.GetRedisObject(ctx)
	//		//ensure timeout is greater than 1s to prevent default timeout and overwrite as needed
	//		redisObj.SetValue("test",p1,1200,true)
	SetValue(ctx context.Context, key string, value interface{}, timeout int, writeIfNotSet bool) error

	//get a value from redis directly into a desired struct
	//	interfaces can be directly sent as reference arguments
	//		eg. usage
	//		type Person struct{
	//			Name string,
	//			Age int,
	//		}
	//		var p1 Person
	//		ctx := context.Background()
	//		redisObj := utils.GetRedisObject(ctx)
	//		redisObj.GetValue("test", &p1)
	GetValue(ctx context.Context, key string, value interface{}) error

	//delete an entry from redis
	DeleteKey(ctx context.Context, key string) error

	//get the TTL for a key set in redis
	GetTTL(ctx context.Context, key string) int

	//check if a key exists in redis
	KeyExists(ctx context.Context, key string) bool

	//add values into redis as a redis hash
	//	accepts a hash key and key-value pairs
	//		eg. usage
	//		type Animal struct{
	//			id		int
	//			name	string
	//			scname 	string
	//			family	string
	//		}
	//		p1 := Person{ name:"elephant", scname:"Loxodonta", family: "mammal"}
	//		dmap := make(map[string]string)
	//		dmap["name"] = p1.name)
	//		dmap["scname"] = p1.scname
	//		dmap["family"] = p1.family
	//		ctx := context.Background()
	//		redisObj := utils.GetRedisObject(ctx)
	//		dkey := fmt.Sprinf("%s_%d", p1.name, p1.id)
	//		redisObj.SetRedisHash(dkey,dmap)
	SetRedisHash(ctx context.Context, key string, kvpairs map[string]string) error

	//fetches individual value for a redis hash key hashmap or all values
	//	all values are fetched if arguments are not provided. only first argument is considered for specific key
	//		eg. usage
	//		type Animal struct{
	//			id		int
	//			name	string
	//			scname 	string
	//			family	string
	//		}
	//		p1 := Person{ name:"elephant", scname:"Loxodonta", family: "mammal"}
	//		dmap := make(map[string]string)
	//		dmap["name"] = p1.name)
	//		dmap["scname"] = p1.scname
	//		dmap["family"] = p1.family
	//		ctx := context.Background()
	//		redisObj := utils.GetRedisObject(ctx)
	//		dkey := fmt.Sprinf("%s_%d", p1.name, p1.id)
	//		dmap, _ := redisObj.GetRedisHashValue(dkey)
	//		dfamily, _:= redisObj.GetRedisHashValue(dkey, "family")
	GetRedisHashValue(ctx context.Context, key string, args ...string) (map[string]string, error)

	DeleteRedisHash(ctx context.Context, key string, fields ...string) error
}

func GetRedisClientItf(url string, port int) RedisClientItf {
	return newRedisClient(url, port)
}
