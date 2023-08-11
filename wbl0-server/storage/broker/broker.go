package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wbl0-server/storage/db"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

var SC *stan.Conn

type OrderRequest struct {
	Order int `json:"order"`
}

func makeOrderJson(oid int) string {
	var or OrderRequest
	or.Order = oid

	jsonString, err := json.Marshal(or)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(jsonString)
}

func Connect() {
	clusterID := loadEnvVar("ClusterID")
	clientID := loadEnvVar("ClientID")
	natsURL := loadEnvVar("NATSUrl")

	sc, err := stan.Connect(
		clusterID,
		clientID,
		stan.NatsURL(natsURL),
		stan.NatsOptions(
			nats.ReconnectWait(time.Second*4),
			nats.Timeout(time.Second*4),
		),
		stan.Pings(5, 3), // Send PINGs every 5 seconds, and fail after 3 PINGs without any response.
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Printf("%s: connection lost, reason: %v", "Broker", reason)
		}),
	)
	checkErr(err)
	SC = &sc
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

func outputErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func SendOrderRequest(oid int) db.Output {
	Connect()
	// Specify the subject you want to subscribe to
	subject := loadEnvVar("NATSSubject")
	log.Println(subject)
	orderJson := makeOrderJson(oid)
	log.Println(orderJson)
	err := (*SC).Publish("order", []byte(orderJson))
	checkErr(err)
	log.Println("test")
	// Subscribe to the subject
	ch := make(chan *stan.Msg)
	sub, err := (*SC).Subscribe("order", func(msg *stan.Msg) {
		ch <- msg
	}, stan.SetManualAckMode(), stan.DeliverAllAvailable())
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	defer (*SC).Close()

	// Capture an interrupt signal to gracefully close the subscription
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	var lastMsg *stan.Msg

	// Start processing messages
	for {
		select {
		case msg := <-ch:
			lastMsg = msg // Store the last received message
			// Process the received message as needed
			fmt.Printf("Received message: %s\n", string(msg.Data))
			msg.Ack()
			var output db.Output
			err := json.Unmarshal([]byte(msg.Data), &output)
			if err != nil {
				log.Println(err)
			}
			if output.OrderUID != "" {
				return output
			}

		case <-signals:
			// Interrupt signal received, exit gracefully
			fmt.Println("Received interrupt signal. Exiting...")
			if lastMsg != nil {
				fmt.Printf("Last received message: %s\n", string(lastMsg.Data))
			}
			var output db.Output
			return output
		}
	}
}
