package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oklog/ulid/v2"
)

// ULIDを生成
func generateID() string {
	return ulid.Make().String()
}

// 会員
type Member struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Address     string    `json:"address" db:"address"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
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

// 貸出情報付き蔵書
type BookWithLend struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Genre     Genre     `json:"genre"`
	Lending   bool      `json:"lending"`
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
	err = db.Get(&key, "SELECT `key` FROM `key` WHERE `id` = (SELECT MAX(`id`) FROM `key`)")
	if err != nil {
		log.Panic(err)
	}

	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		log.Panic(err)
	}

	e := echo.New()
	e.Debug = true
	e.Use(middleware.Logger())

	api := e.Group("/api")
	{
		api.POST("/initialize", initializeHandler)

		membersAPI := api.Group("/members")
		{
			membersAPI.POST("", postMemberHandler)
			membersAPI.GET("", getMembersHandler)
			membersAPI.GET("/:id", getMemberHandler)
			membersAPI.PATCH("/:id", patchMemberHandler)
			membersAPI.DELETE("/:id", banMemberHandler)
			membersAPI.GET("/:id/qrcode", getMemberQRCodeHandler)
		}

		booksAPI:= api.Group("/books")
		{
			booksAPI.GET(":id", getBookHandler)
		}
	}

	e.Logger.Fatal(e.Start(":8080"))
}

type InitializeHandlerRequest struct {
	Key string `json:"key"`
}

type InitializeHandlerResponse struct {
	Language string `json:"language"`
}

/*
---------------------------------------------------------------
Initialization API
---------------------------------------------------------------
*/

// 初期化用ハンドラ
func initializeHandler(c echo.Context) error {
	var req InitializeHandlerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(req.Key) != 16 {
		return echo.NewHTTPError(http.StatusBadRequest, "key must be 16 characters")
	}

	cmd := exec.Command("sh", "../sql/init_db.sh")
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, err = db.Exec("INSERT INTO `key` (`key`) VALUES (?)", req.Key)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	block, err = aes.NewCipher([]byte(req.Key))
	if err != nil {
		log.Panic(err.Error())
	}

	return c.JSON(http.StatusOK, InitializeHandlerResponse{
		Language: "Go",
	})
}

/*
---------------------------------------------------------------
Members API
---------------------------------------------------------------
*/

type PostMemberRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

// 会員登録
func postMemberHandler(c echo.Context) error {
	var req PostMemberRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if req.Name == "" || req.Address == "" || req.PhoneNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name, address, phoneNumber are required")
	}

	id := generateID()

	_, err := db.Exec("INSERT INTO `member` (`id`, `name`, `address`, `phone_number`, `banned`, `created_at`) VALUES (?, ?, ?, ?, false, ?)",
		id, req.Name, req.Address, req.PhoneNumber, time.Now())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var res Member
	err = db.Get(&res, "SELECT * FROM `member` WHERE `id` = ?", id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

const memberPageLimit = 500

type GetMembersResponse struct {
	Members []Member `json:"members"`
	Total   int      `json:"total"`
}

// 会員一覧を取得 (ページネーションあり)
func getMembersHandler(c echo.Context) error {
	pageStr := c.QueryParam("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 前ページの最後の会員ID
	// シーク法をフロントエンドでは実装したが、バックエンドは力尽きた
	_ = c.QueryParam("last_member_id")

	members := []Member{}
	err = db.Select(&members, "SELECT * FROM `member` WHERE `banned` = false ORDER BY `id` LIMIT ? OFFSET ?",
		memberPageLimit, (page-1)*memberPageLimit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if len(members) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "no members to show in this page")
	}

	var total int
	err = db.Get(&total, "SELECT COUNT(*) FROM `member`")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, GetMembersResponse{
		Members: members,
		Total:   total,
	})
}

// 会員を取得
func getMemberHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	encrypted := c.QueryParam("encrypted")
	if encrypted == "true" {
		var err error
		id, err = decrypt(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	} else if encrypted != "" && encrypted != "false" {
		return echo.NewHTTPError(http.StatusBadRequest, "encrypted must be boolean value")
	}

	member := Member{}
	err := db.Get(&member, "SELECT * FROM `member` WHERE `id` = ? AND `banned` = false", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, member)
}

type PatchMemberRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

// 会員情報編集
func patchMemberHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	var req PatchMemberRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if req.Name == "" && req.Address == "" && req.PhoneNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name, address or phoneNumber is required")
	}

	// 会員の存在を確認
	err := db.Get(&Member{}, "SELECT * FROM `member` WHERE `id` = ? AND `banned` = false", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	query := "UPDATE `member` SET "
	params := []any{}
	if req.Name != "" {
		query += "`name` = ?, "
		params = append(params, req.Name)
	}
	if req.Address != "" {
		query += "`address` = ?, "
		params = append(params, req.Address)
	}
	if req.PhoneNumber != "" {
		query += "`phone_number` = ?, "
		params = append(params, req.PhoneNumber)
	}
	query = query[:len(query)-2] // 最後の", "を削除
	query += " WHERE `id` = ?"
	params = append(params, id)

	_, err = db.Exec(query, params...)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// 会員をBAN
func banMemberHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	// 会員の存在を確認
	err := db.Get(&Member{}, "SELECT * FROM `member` WHERE `id` = ? AND `banned` = false", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, err = db.Exec("UPDATE `member` SET `banned` = true WHERE `id` = ?", id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// 会員証用のQRコードを取得
func getMemberQRCodeHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	// 会員の存在確認
	err := db.Get(&Member{}, "SELECT * FROM `member` WHERE `id` = ? AND `banned` = false", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	qrFileLock.Lock()
	defer qrFileLock.Unlock()

	qrCode, err := generateQRCode(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Blob(http.StatusOK, "image/png", qrCode)
}

/*
---------------------------------------------------------------
Books API
---------------------------------------------------------------
*/

// 蔵書を取得
func getBookHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	encrypted := c.QueryParam("encrypted")
	if encrypted == "true" {
		var err error
		id, err = decrypt(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	} else if encrypted != "" && encrypted != "false" {
		return echo.NewHTTPError(http.StatusBadRequest, "encrypted must be boolean value")
	}

	book := Book{}
	err := db.Get(&book, "SELECT * FROM `book` WHERE `id` = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, book)
}

/*
---------------------------------------------------------------
Lending API
---------------------------------------------------------------
*/
