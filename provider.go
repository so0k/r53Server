package main

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
)

type cloud interface {
	List(ctx context.Context) ([]zone, error)
	ZonesString() string
}

func newProvider(provider, awsAccessKey, awsSecretKey string, zoneIds []string) (cloud, error) {
	conf := initAwsConfig(awsAccessKey, awsSecretKey)

	p := r53Provider{
		client:  route53.New(session.New(conf)),
		zoneIds: zoneIds,
	}
	logrus.Info(zoneIds)
	return &p, nil
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
