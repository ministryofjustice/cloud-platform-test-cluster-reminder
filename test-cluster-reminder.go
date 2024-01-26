package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
)

type SlackRequestBody struct {
	Text string `json:"text"`
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
	})

	if err != nil {
		log.Println("Error creating session,", err)
		return
	}

	svc := eks.New(sess)

	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		log.Println("Error: SLACK_WEBHOOK_URL environment variable not set")
		return
	}

	result, err := svc.ListClusters(nil)
	if err != nil {
		log.Println("Error listing clusters,", err)
		return
	}

	for i := len(result.Clusters) - 1; i >= 0; i-- {
		if strings.Contains(*result.Clusters[i], "live") || strings.Contains(*result.Clusters[i], "manager") {
			result.Clusters = append(result.Clusters[:i], result.Clusters[i+1:]...)
		}
	}

	var openEmoji = ":donut_spin: :donut_spin: :friday_yayday:"
	var closeEmoji = ":friday_yayday: :donut_spin: :donut_spin:"

	var clustersString string
	for _, cluster := range result.Clusters {
		clustersString += *cluster + "\n"
	}

	if clustersString == "" {
		clustersString = openEmoji + " *Yay, there are no test clusters to cleanup!*  " + closeEmoji
	} else {
		clustersString = openEmoji + " *Friday test cluster cleanup reminder!* " + closeEmoji + " \n \n Please delete your test clusters before signing off for the weekend! \n ```" + clustersString + "```"
	}

	err = SendSlackNotification(webhookURL, clustersString)
	if err != nil {
		log.Println("Error sending slack notification,", err)
	}
}

func SendSlackNotification(webhookUrl string, msg string) error {
	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("non-ok response returned from slack")
	}
	return nil
}
