package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

var (
	ctx = context.Background()
)

func main() {
	zoneID := "ZONEID"
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ACCESS_KEY", "SECRET_KEY", "")))
	if err != nil {
		fmt.Printf("Error creating config: %s", err)
	}

	client := route53.NewFromConfig(cfg)
	//listHostedZones(client)
	//listRecordSets(client, zoneID, 10)
	changeRecordSets(client, "mktest-api.arun-vpn.internal", "10.0.1.101", zoneID, types.ChangeActionUpsert)
}
func listHostedZones(client *route53.Client) {
	output, err := client.ListHostedZones(ctx, &route53.ListHostedZonesInput{})
	if err != nil {
		fmt.Printf("Error listing zones: %s", err)
	}

	for _, h := range output.HostedZones {
		log.Printf("id=%s name=%s", *h.Id, *h.Name)
	}
}

func listRecordSets(client *route53.Client, zoneID string, maxItems int32) {
	output, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
		HostedZoneId: &zoneID,
		MaxItems:     &maxItems,
	})
	if err != nil {
		fmt.Printf("Error listing record sets: %s", err)
	}

	for _, r := range output.ResourceRecordSets {
		log.Printf("name=%s type=%s", *r.Name, r.Type)
	}
}

func changeRecordSets(client *route53.Client, name, value, zoneID string, action types.ChangeAction) {
	_, err := client.ChangeResourceRecordSets(ctx, generateChangeResourceRecordSetsInput(name, value, zoneID, action))
	if err != nil {
		fmt.Printf("Error in changeRecordSets: %s", err)
	}
}

func generateChangeResourceRecordSetsInput(name, value, zoneID string, action types.ChangeAction) *route53.ChangeResourceRecordSetsInput {
	r := &types.ResourceRecordSet{
		Name: aws.String(name),
		Type: types.RRTypeA,
		TTL:  aws.Int64(60),
		ResourceRecords: []types.ResourceRecord{
			types.ResourceRecord{
				Value: aws.String(value),
			},
		},
	}
	return &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action:            action,
					ResourceRecordSet: r,
				},
			},
		},
	}
}
