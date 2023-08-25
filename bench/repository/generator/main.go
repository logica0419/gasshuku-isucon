package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/repository"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

func getEnvOrDefault(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	return defaultValue
}

func main() {
	host := getEnvOrDefault("DB_HOST", "127.0.0.1")
	port := getEnvOrDefault("DB_PORT", "3306")
	user := getEnvOrDefault("DB_USER", "isucon")
	pass := getEnvOrDefault("DB_PASS", "isucon")
	name := getEnvOrDefault("DB_NAME", "isulibrary")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo", user, pass, host, port, name)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	tx, err := db.Beginx()
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	initData := repository.InitData{
		Members: []*model.MemberWithLending{},
		Books:   []*model.BookWithLending{},
	}

	for i := 0; i < 76; i++ {
		members := []*model.MemberWithLending{}
		for i := 0; i < 992; i++ {
			member := model.NewMember()
			members = append(members, member)
			time.Sleep(time.Duration((1 + rand.Intn(100))) * time.Microsecond)
		}
		_, err = tx.NamedExec(
			"INSERT INTO member (`id`, `name`, `address`, `phone_number`, `banned`, `created_at`) "+
				"VALUES (:id, :name, :address, :phone_number, :banned, :created_at)",
			members,
		)
		if err != nil {
			log.Panic(err)
		}
	}
	err = tx.Select(&initData.Members, "SELECT * FROM member")
	if err != nil {
		log.Panic(err)
	}

	for i := 0; i < 127; i++ {
		books := []*model.BookWithLending{}
		for i := 0; i < 866; i++ {
			book := model.NewBook()
			books = append(books, book)
			time.Sleep(time.Duration((1 + rand.Intn(100))) * time.Microsecond)
		}
		_, err = tx.NamedExec(
			"INSERT INTO book (`id`, `title`, `author`, `genre`, `created_at`) "+
				"VALUES (:id, :title, :author, :genre, :created_at)",
			books,
		)
		if err != nil {
			log.Panic(err)
		}
	}
	err = tx.Select(&initData.Books, "SELECT * FROM book")
	if err != nil {
		log.Panic(err)
	}

	_, err = tx.Exec("INSERT INTO `key` (`key`) VALUES (?)", utils.RandStringWithSign(16))
	if err != nil {
		log.Panic(err)
	}

	_ = tx.Commit()

	f, err := os.Create("../init_data.json")
	if err != nil {
		log.Panic(err)
	}
	r, err := utils.EncodeJson(initData)
	if err != nil {
		log.Panic(err)
	}
	_, err = io.Copy(f, r)
	if err != nil {
		log.Panic(err)
	}
}
