package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const PushEndpoint = "https://push.statuscake.com/"

func sendStatus() error {

	pk := os.Getenv("STATUS_CAKE_PK")
	test_id := os.Getenv("STATUS_CAKE_TEST_ID")

	t0 := time.Now().Unix()
	url := PushEndpoint + fmt.Sprintf("?PK=%v&TestID=%v&time=%v", pk, test_id, t0)
	log.Println("Sending status", url)
	_, err := http.Get(url)
	if err != nil {
		log.Println("Send status failed")
		return err
	}
	return err
}

func main() {
	err := sendStatus()
	if err != nil {
		log.Fatalln(err)
	}
}
