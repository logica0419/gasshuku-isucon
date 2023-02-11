package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

// ULIDを生成
func generateID() string {
	return ulid.Make().String()
}

// 図書館
type Library struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Address     string `json:"address" db:"address"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
}

// 会員
type Member struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Address     string    `json:"address" db:"address"`
	PhoneNumber string    `json:"phoneNumber" db:"phone_number"`
	LibraryID   string    `json:"libraryId" db:"library_id"`
	Banned      bool      `json:"banned" db:"banned"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// 図書分類
type Genre int

// 国際十進分類法に従った図書分類
const (
	General         Genre = iota // 総記
	Philosophy                   // 哲学・心理学
	Religion                     // 宗教・神学
	SocialScience                // 社会科学
	Vacant                       // 未定義
	Mathematics                  // 数学・自然科学
	AppliedSciences              // 応用科学・医学・工学
	Arts                         // 芸術
	Literature                   // 言語・文学
	Geography                    // 地理・歴史
)

// 蔵書
type Book struct {
	ID        string    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Author    string    `json:"author" db:"author"`
	Genre     Genre     `json:"genre" db:"genre"`
	LibraryID string    `json:"library_id" db:"library_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// 貸出記録
type Lending struct {
	ID        string    `json:"id" db:"id"`
	MemberID  string    `json:"member_id" db:"member_id"`
	BookID    string    `json:"book_id" db:"book_id"`
	Due       time.Time `json:"due" db:"due"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// 蔵書取り寄せリクエスト
type Order struct {
	ID        string    `json:"id" db:"id"`
	BookID    string    `json:"book_id" db:"book_id"`
	FromID    string    `json:"from_id" db:"from_id"`
	ToID      string    `json:"to_id" db:"to_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// 貸出情報付き蔵書
type BookWithLend struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Genre     Genre     `json:"genre"`
	Lending   bool      `json:"lending"`
	LibraryID string    `json:"library_id"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sqlx.DB

func getEnvOrDefault(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	return defaultValue
}

var (
	block      cipher.Block
	qrFileLock sync.Mutex
)

// AES + CTRモード + base64エンコードでテキストを暗号化
func encrypt(plainText string) (string, error) {
	cipherText := make([]byte, aes.BlockSize+len([]byte(plainText)))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AES + CTRモード + base64エンコードで暗号化されたテキストを複合
func decrypt(cipherText string) (string, error) {
	cipherByte, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	decryptedText := make([]byte, len([]byte(cipherByte[aes.BlockSize:])))
	decryptStream := cipher.NewCTR(block, []byte(cipherByte[:aes.BlockSize]))
	decryptStream.XORKeyStream(decryptedText, []byte(cipherByte[aes.BlockSize:]))
	return string(decryptedText), nil
}

const qrCodeFileName = "../images/qr.png"

// QRコードを生成
func generateQRCode(id string) ([]byte, error) {
	encryptedID, err := encrypt(id)
	if err != nil {
		return nil, err
	}

	/*
		生成するQRコードの仕様
		 - PNGフォーマット
		 - QRコードの1モジュールは1ピクセルで表現
		 - バージョン5 (37x37ピクセル、マージン含め45x45ピクセル)
		 - エラー訂正レベルM (15%)
	*/
	_, err = exec.
		Command("qrencode", "-o", qrCodeFileName, "-t", "PNG", "-s", "1", "-v", "5", "--strict-version", "-l", "M", encryptedID).
		Output()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(qrCodeFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func main() {
	host := getEnvOrDefault("DB_HOST", "127.0.0.1")
	port := getEnvOrDefault("DB_PORT", "3306")
	user := getEnvOrDefault("DB_USER", "isucon")
	pass := getEnvOrDefault("DB_PASS", "isucon")
	name := getEnvOrDefault("DB_NAME", "isulibrary")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo", user, pass, host, port, name)

	var err error
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Panic(err)
	}

	var key string
	err = db.Get(&key, "SELECT `key` FROM `key` WHERE `id` = 1")
	if err != nil {
		log.Panic(err)
	}

	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		log.Panic(err)
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
