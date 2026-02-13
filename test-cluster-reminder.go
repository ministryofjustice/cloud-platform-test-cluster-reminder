package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/ministryofjutice/cloud-platform-test-cluster-reminder/handlers"
	"github.com/ministryofjutice/cloud-platform-test-cluster-reminder/utils"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))

	if err != nil {
		log.Println("Error loading AWS config,", err)
		return
	}

	svc := eks.NewFromConfig(cfg)
	accountName, err := handlers.GetAWSAccountName(ctx, cfg)
	if err != nil {
		log.Println("Error getting AWS account name,", err)
		accountName = "unknown-account"
	}

	// webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	// if webhookURL == "" {
	// 	log.Println("Error: SLACK_WEBHOOK_URL environment variable not set")
	// 	return
	// }

	// slackChannel := os.Getenv("SLACK_CHANNEL")
	// if slackChannel == "" {
	// 	log.Println("Error: SLACK_CHANNEL environment variable not set")
	// 	return
	// }

	result, err := svc.ListClusters(ctx, &eks.ListClustersInput{})
	if err != nil {
		log.Println("Error listing clusters,", err)
		return
	}

	accountClusters := []string{"cloud-platform-development", "cloud-platform-preproduction", "cloud-platform-nonlive", "cloud-platform-live"}

	for i := len(result.Clusters) - 1; i >= 0; i-- {
		if utils.HasAnyPrefix(result.Clusters[i], accountClusters) {
			result.Clusters = append(result.Clusters[:i], result.Clusters[i+1:]...)
		}
	}

	var openEmoji = ":donut_spin: :donut_spin: :friday_yayday:"
	var closeEmoji = ":friday_yayday: :donut_spin: :donut_spin:"

	var clustersString string
	for _, cluster := range result.Clusters {
		clustersString += cluster + "\n"
	}

	if clustersString == "" {
		clustersString = openEmoji + " *Yay, there are no test clusters to cleanup!*  " + closeEmoji
	} else {
		clustersString = openEmoji + " *Friday test cluster cleanup reminder!* " + closeEmoji + " \n \n Please delete your test clusters before signing off for the weekend! \n ```" + accountName + "\n" + clustersString + "```"
	}

	log.Println(clustersString)

	// err = handlers.SendSlackNotification(webhookURL, clustersString, slackChannel)
	// if err != nil {
	// 	log.Println("Error sending slack notification,", err)
	// }
}
