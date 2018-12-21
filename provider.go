package main

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type cloud interface {
	List(ctx context.Context) ([]zone, error)
	ZonesString() string
}

func newProviders(provider, awsAccessKey, awsSecretKey string, config *Config) ([]cloud, error) {
	awsConf := initAwsConfig(awsAccessKey, awsSecretKey)

	var providers []cloud
	for _, r := range config.Roles {
		sess := session.New(awsConf)
		var p *r53Provider
		if r.RoleArn == "none" {
			// use AWS Chain credentials
			p = &r53Provider{
				client:  route53.New(sess),
				zoneIds: r.Zones,
			}
		} else {
			// use AWS Chain with AssumeRole credential provider
			creds := stscreds.NewCredentials(sess, r.RoleArn)
			p = &r53Provider{
				client:  route53.New(sess, &aws.Config{Credentials: creds}),
				zoneIds: r.Zones,
			}
		}
		providers = append(providers, p)
	}
	return providers, nil
}

func initAwsConfig(accessKey, secretKey string) *aws.Config {
	awsConfig := aws.NewConfig()
	creds := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			},
		},
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{},
	})
	awsConfig.WithCredentials(creds)
	return awsConfig
}
