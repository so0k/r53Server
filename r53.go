package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

type r53Provider struct {
	zoneIds []string
	client  *route53.Route53
	ctx     context.Context
}

func (c *r53Provider) ZonesString() string {
	return strings.Join(c.zoneIds, ",")
}

// List returns the A records from all provided zoneIds.
func (c *r53Provider) List(ctx context.Context) (zones []zone, err error) {
	for _, zoneId := range c.zoneIds {
		result, err := c.client.GetHostedZone(&route53.GetHostedZoneInput{
			Id: aws.String(zoneId),
		})
		if err != nil {
			panic(fmt.Sprintf("failed to find zone, %s, %v", zoneId, err))
		}
		z := zone{
			Name: aws.StringValue(result.HostedZone.Name),
		}

		err = c.client.ListResourceRecordSetsPagesWithContext(ctx, &route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String(zoneId),
		}, func(p *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			for _, rrset := range p.ResourceRecordSets {
				if *rrset.Type == "A" {
					z.Records = append(z.Records, strings.TrimSuffix(aws.StringValue(rrset.Name), "."))
				}
			}
			return true // continue paging
		})
		zones = append(zones, z)
		if err != nil {
			panic(fmt.Sprintf("failed to list records for zone, %s, %v", zoneId, err))
		}
	}

	return zones, nil
}
