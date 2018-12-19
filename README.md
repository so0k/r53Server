# r53Server

[![Docker Repository on Quay](https://quay.io/repository/swo0k/r53server/status?token=6ae97470-ce13-47d4-a054-7b992f6507b2 "Docker Repository on Quay")](https://quay.io/repository/swo0k/r53server)
[![Build Status](https://travis-ci.org/so0k/r53Server.svg?branch=master)](https://travis-ci.org/so0k/r53Server)

Static server for A records in r53 zones.

based on [jessfraz/s3server](https://github.com/jessfraz/s3server)

## Usage

```bash
Server to index & view recods in r53 zones.
 Version: v0.1.1
 Build: 9089a07
  -aws-access-key-id string
        AWS access key
  -aws-secret-access-key string
        AWS access secret
  -interval string
        interval to generate new index.html's at (default "5m")
  -p string
        port for server to run on (default "8080")
  -v    print version and exit (shorthand)
  -version
        print version and exit
  -zone value
        Route53 Zone Id to fetch records from (can be repeated)
```

## Docker

```bash
docker run -d \
  -p 8080:8080 \
  -e AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY \
  --tmpfs /tmp \
  quay.io/swo0k/r53server:v0.1.1 \
    -zone Z2UE.......... \
    -zone Z1W...........
```
