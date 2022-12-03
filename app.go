package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Author struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {

	fmt.Println("Hello World")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	pong, err := client.Ping(client.Context()).Result()
	fmt.Println(pong, err)

	if err != nil {
		log.Fatalf("Error connecting to redis")
	}

	author := Author{Name: "Ben", Age: 85}
	create(client, "id4", author)
	fmt.Println(read(client, "id4"))

	author1 := Author{Name: "Jhon", Age: 55}
	create(client, "id5", author1)
	fmt.Println(read(client, "id5"))

	addItemOnLis(client, "Google", "Facebook", "Amazon")
	fmt.Println(client.LRange(client.Context(), "companies", 0, -1).Val())

	readAllKeysOfAuthor(client)

	author2 := Author{Name: "Jhon", Age: 99}
	updateAuthor(client, "id5", author2)
	fmt.Println(read(client, "id5"))

	deleteAuthor(client, "id5")
	fmt.Println(read(client, "id5"))
}

func create(client *redis.Client, id string, author Author) {
	dataJson, err := json.Marshal(author)
	if err != nil {
		fmt.Println(err)
	}

	err = client.Set(client.Context(), id, dataJson, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func read(client *redis.Client, id string) Author {
	val, err := client.Get(client.Context(), id).Result()
	if err != nil {
		fmt.Println(fmt.Sprint("Error reading key: ", id))
	}

	var author Author
	err = json.Unmarshal([]byte(val), &author)
	if err != nil {
		fmt.Println(err)
	}

	return author
}

func addItemOnLis(client *redis.Client, companies ...string) {
	err := client.LPush(client.Context(), "companies", companies).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func readAllKeysOfAuthor(client *redis.Client) {
	keys, err := client.Keys(client.Context(), "id*").Result()
	if err != nil {
		fmt.Println(err)
	}

	for _, key := range keys {
		fmt.Println(fmt.Sprint(key, ":", read(client, key)))
	}

}

func updateAuthor(client *redis.Client, id string, author Author) {
	dataJson, err := json.Marshal(author)
	if err != nil {
		fmt.Println(err)
	}

	err = client.Set(client.Context(), id, dataJson, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func deleteAuthor(client *redis.Client, id string) {
	err := client.Del(client.Context(), id).Err()
	if err != nil {
		fmt.Println(err)
	}
}
