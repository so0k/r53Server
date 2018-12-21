# r53Server

[![Docker Repository on Quay](https://quay.io/repository/swo0k/r53server/status?token=6ae97470-ce13-47d4-a054-7b992f6507b2 "Docker Repository on Quay")](https://quay.io/repository/swo0k/r53server)
[![Build Status](https://travis-ci.org/so0k/r53Server.svg?branch=master)](https://travis-ci.org/so0k/r53Server)

Static server for A records in r53 zones.

based on [jessfraz/s3server](https://github.com/jessfraz/s3server)

## Usage

```bash
$ ./r53Server -h

 Server to index & view records in r53 zones.
 Version: v0.2.0
 Build: 6ce6dfd
  -aws-access-key-id string
        AWS access key
  -aws-secret-access-key string
        AWS access secret
  -config string
        config file (default "config.yaml")
  -interval string
        interval to generate new index.html's at (default "5m")
  -p string
        port for server to run on (default "8080")
  -v    print version and exit (shorthand)
  -version
        print version and exit
```

`config.yaml` format

```yaml
# roleArn `none` is special keyword
roles:
- roleArn: none
  zones: []
  - Z...
- roleArn: "arn:aws:iam::..."
  zones:
  - Z...
```

## Docker

```bash
docker run -d \
  -p 8080:8080 \
  -e AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY \
  -v $PWD/config.yaml:/config.yaml \
  --tmpfs /tmp \
  quay.io/swo0k/r53server:v0.2.0 \
    -p 8080
```
