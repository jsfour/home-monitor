package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

const PublicIpEndpoint = "https://api.ipify.org"
const PushEndpoint = "https://push.statuscake.com/"

func getPublicIp() (string, error) {
	res, err := http.Get(PublicIpEndpoint)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	s := buf.String()

	return s, nil
}

func sendStatus() error {

	pk := os.Getenv("STATUS_CAKE_PK")
	test_id := os.Getenv("STATUS_CAKE_TEST_ID")

	t0 := 0
	url := PushEndpoint + fmt.Sprintf("?PK=%v&TestID=%v&time=%v", pk, test_id, t0)
	log.Println("Sending status", url)
	res, err := http.Get(url)
	if err != nil {
		log.Println("Send status failed")
		return err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	s := buf.String()
	if s != "success" {
		log.Fatalln(s)
	}
	return err
}

func main() {
	err := sendStatus()
	if err != nil {
		log.Fatalln(err)
	}
}
