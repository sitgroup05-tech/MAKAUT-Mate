package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		log.Fatal("Error loading .env file")
	}

	db_string := buildDSNString(
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_DB_NAME"),
		os.Getenv("PG_USER_NAME"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_SSL_MODE"),
	)
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if host == "" || port == "" {
		log.Fatalf("Host or Port Can't Be NULL Given Host:%s PORT:%s\n", host, port)
	}

	rpc_host := os.Getenv("RPC_HOST")
	rpc_port := os.Getenv("RPC_PORT")

	if rpc_host == "" || rpc_port == "" {
		log.Fatalf("RPC Host or Port Can't be NULL Given RPC_HOST:%s RPC_PORT:%s\n", rpc_host, rpc_port)
	}

	secret_key := os.Getenv("SECRET_KEY")

	if secret_key == "" {
		log.Fatalln("SECRET KEY can't be NULL")
	}

	conf := Config{
		host: host,
		port: port,
		db: DBConfig{
			dbn: db_string,
		},
		rpc: RPCConfig{
			host: rpc_host,
			port: rpc_port,
		},
		secretKey: secret_key,
	}

	app := NewApplication(conf)
	defer app.Close()

	app.run(app.mound())
}

func buildDSNString(
	host, port, db,
	usr, pwd,
	ssl string) string {

	if host == "" ||
		port == "" ||
		db == "" ||
		usr == "" ||
		pwd == "" ||
		ssl == "" {
		log.Fatalf(
			"Invalid DB_DATA Given host=%s port%s db=%s user=%s password=%s sslmode=%s\n",
			host, port, usr, pwd, db, ssl,
		)
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, usr, pwd, db, ssl,
	)
}
