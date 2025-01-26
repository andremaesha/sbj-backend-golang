package bootstrap

import "fmt"

func dsn(host, port, user, password, dbName string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host,
		user,
		password,
		dbName,
		port,
	)
}
