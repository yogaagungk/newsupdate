package config

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

//InitialRedisConn , open connection dan konfigurasi pool yang digunakan
//Cache menggunakan Redis
func InitialRedisConn() redis.Conn {
	redisPool := &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}

			if _, err := c.Do("AUTH", "12345"); err != nil {
				c.Close()
				return nil, err
			}

			return c, err
		},
	}

	redisConn := redisPool.Get()

	ping(redisConn)

	return redisConn
}

func ping(c redis.Conn) error {
	// Send PING command to Redis
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}

	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(pong, err)
	if err != nil {
		return err
	}

	fmt.Printf("PING Response = %s\n", s)
	// Output: PONG

	return nil
}
