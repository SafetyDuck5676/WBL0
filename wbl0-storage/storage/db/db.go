package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type Order struct {
	OrderUID          string   `json:"order_uid"`
	Entry             string   `json:"entry"`
	InternalSignature string   `json:"internal_signature"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Items  `json:"items"`
	Locale            string   `json:"locale"`
	CustomerID        string   `json:"customer_id"`
	TrackNumber       string   `json:"track_number"`
	DeliveryService   string   `json:"delivery_service"`
	Shardkey          string   `json:"shardkey"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          int      `json:"oof_shard"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}
type Items struct {
	ChrtID     int    `json:"chrt_id"`
	Price      int    `json:"price"`
	Rid        string `json:"rid"`
	Name       string `json:"name"`
	Sale       int    `json:"sale"`
	Size       string `json:"size"`
	TotalPrice int    `json:"total_price"`
	NmID       int    `json:"nm_id"`
	Brand      string `json:"brand"`
	Status     int    `json:"status"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"price"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Output struct {
	OrderID         int    `json:"order"`
	OrderUID        string `json:"order_uid"`
	Entry           string `json:"entry"`
	CustomerID      string `json:"customer_id"`
	TrackNumber     string `json:"track_number"`
	DeliveryService string `json:"delivery_service"`
}

type OrderRequest struct {
	Order int `json:"order"`
}

func ConnectDB() {
	var err error
	host := loadEnvVar("DBhost")
	port := loadEnvVar("DBport")
	user := loadEnvVar("DBuser")
	password := loadEnvVar("DBpassword")
	dbname := loadEnvVar("DBname")

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	DB, err = sql.Open("postgres", psqlconn)
	checkErr(err)

	err = DB.Ping()
	checkErr(err)
}

func GetOrderById(oid int) Order {
	ConnectDB()
	var o Order

	sql := "SELECT " +
		"order_uid, " +
		"track_number, " +
		"entry, " +
		"locale, " +
		"internal_signature, " +
		"customer_id, " +
		"delivery_service, " +
		"shardkey, " +
		"sm_id, " +
		"date_created, " +
		"oof_shard, " +
		"transaction, " +
		"request_id, " +
		"currency, " +
		"provider, " +
		"amount, " +
		"payment_dt, " +
		"bank, " +
		"delivery_cost, " +
		"goods_total, " +
		"custom_fee, " +
		"deliveries.name, " +
		"phone, " +
		"zip, " +
		"city, " +
		"address, " +
		"region, " +
		"email, " +
		"chrt_id, " +
		"price, " +
		"rid, " +
		"items.name AS itemname, " +
		"sale, " +
		"size, " +
		"total_price, " +
		"nm_id, " +
		"brand, " +
		"status " +
		"FROM orders " +
		"INNER JOIN payments ON orders.payment_id = payments.id " +
		"INNER JOIN deliveries ON orders.delivery_id = deliveries.id " +
		"INNER JOIN item_order ON orders.id = item_order.order_id " +
		"INNER JOIN items ON item_order.item_id = items.id " +
		"WHERE orders.id = $1"
	rows, err := DB.Query(sql, oid)
	defer rows.Close()

	checkErr(err)

	for rows.Next() {
		var item Items
		err = rows.Scan(
			&o.OrderUID,
			&o.TrackNumber,
			&o.Entry,
			&o.Locale,
			&o.InternalSignature,
			&o.CustomerID,
			&o.DeliveryService,
			&o.Shardkey,
			&o.SmID,
			&o.DateCreated,
			&o.OofShard,
			&o.Payment.Transaction,
			&o.Payment.RequestId,
			&o.Payment.Currency,
			&o.Payment.Provider,
			&o.Payment.Amount,
			&o.Payment.PaymentDt,
			&o.Payment.Bank,
			&o.Payment.DeliveryCost,
			&o.Payment.GoodsTotal,
			&o.Payment.CustomFee,
			&o.Delivery.Name,
			&o.Delivery.Phone,
			&o.Delivery.Zip,
			&o.Delivery.City,
			&o.Delivery.Address,
			&o.Delivery.Region,
			&o.Delivery.Email,
			&item.ChrtID,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status)
		o.Items = append(o.Items, item)
		checkErr(err)
	}
	return o
}

func AddOrder(o Order) {
	var payment_id int
	var delivery_id int
	var order_id int

	insertPayment := `INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`
	err := DB.QueryRow(insertPayment, o.Payment.Transaction, o.Payment.RequestId, o.Payment.Currency, o.Payment.Provider, o.Payment.Amount, o.Payment.PaymentDt, o.Payment.Bank, o.Payment.DeliveryCost, o.Payment.GoodsTotal, o.Payment.CustomFee).Scan(&payment_id)
	checkErr(err)

	insertDelivery := `INSERT INTO "deliveries" ("name", "phone", "zip", "city", "address", "region", "email") VALUES ($1, $2, $3, $4, $5, $6, $7);`
	err = DB.QueryRow(insertDelivery, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip, o.Delivery.City, o.Delivery.Address, o.Delivery.Region, o.Delivery.Email).Scan(&delivery_id)
	checkErr(err)

	insertOrder := `INSERT INTO "orders" ("order_uid", "track_number", "entry", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "oof_shard", "payment_id", "delivery_id") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`
	err = DB.QueryRow(insertOrder, o.OrderUID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature, o.CustomerID, o.DeliveryService, o.Shardkey, o.SmID, o.OofShard, payment_id, delivery_id).Scan(&order_id)
	checkErr(err)

	insertItems := "INSERT INTO items (chrt_id, rid, name, sale, size, nm_id, brand, status, total_price, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);"
	for _, value := range o.Items {
		var lastItemId int
		err := DB.QueryRow(insertItems, value.ChrtID, value.Rid, value.Name, value.Sale, value.Size, value.NmID, value.Brand, value.Status, value.TotalPrice, value.Price).Scan(&lastItemId)
		checkErr(err)

		insertItemOrder := `INSERT INTO "item_order" ("item_id", "order_id") VALUES ($1, $2);`
		_, err = DB.Exec(insertItemOrder, lastItemId, order_id)
	}
}

func loadEnvVar(envVar string) string {
	err := godotenv.Load(".env")
	checkErr(err)
	return os.Getenv(envVar)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
