package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
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

func sendDns() error {
	log.Println("Updating dns")
	sess := session.Must(session.NewSession(&aws.Config{}))

	svc := route53.New(sess)

	ip, err := getPublicIp()

	if err != nil {
		return err
	}

	dnsName := os.Getenv("DNS_HOSTNAME")
	zoneID := os.Getenv("ROUTE53_ZONE_ID")
	var ttl int64 = 300

	recordSetQuery := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneID),
		StartRecordName: aws.String(dnsName),
	}

	res, err := svc.ListResourceRecordSets(recordSetQuery)
	if err != nil {
		return err
	}

	var changes []*route53.Change
	for _, recSet := range res.ResourceRecordSets {
		if *recSet.Name != dnsName {
			continue
		}
		for _, record := range recSet.ResourceRecords {
			log.Println(*record.Value)
			if *record.Value == ip {
				log.Println("No change in IP")
				return nil
			}
			// need to change

			ttl = *recSet.TTL
			chg := &route53.Change{
				Action:            aws.String("DELETE"),
				ResourceRecordSet: recSet,
			}
			changes = append(changes, chg)

			break
		}
	}

	newChange := &route53.Change{
		Action: aws.String("CREATE"),
		ResourceRecordSet: &route53.ResourceRecordSet{
			Name: aws.String(dnsName),
			Type: aws.String("A"),
			TTL:  aws.Int64(ttl),
			ResourceRecords: []*route53.ResourceRecord{
				{
					Value: aws.String(ip),
				},
			},
		},
	}

	changes = append(changes, newChange)

	changeSet := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
		},
		HostedZoneId: aws.String(zoneID),
	}

	_, err = svc.ChangeResourceRecordSets(changeSet)
	if err != nil {
		return err
	}

	return nil
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
	err = sendDns()
	if err != nil {
		log.Fatalln(err)
	}
}
