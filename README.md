![Go](https://github.com/mkromkamp/DnsChecker/workflows/Go/badge.svg)

# DnsChecker

Small application to query specific dns servers for host names.

## Usage

In order to run DnsChecker you have three options;

### Docker

Run the [docker image](https://hub.docker.com/r/mkromkamp/dnschecker) from docker hub;

``` bash
docker run -p 8080:8080 mkromkamp/dnschecker:latest
```

### Binary

Download the binary from the releases page;

``` bash
chmod +x dnschecker && ./dnschecker
```

## Configuration

| Var | Default | Description |
|---|---|---|
| APP_NAME | DnsChecker-Local | Application name, mainly used for logging |
| APP_NAMESERVER | 8.8.8.8 | Name server to query; seperate by comma for multiple |
| APP_HOSTNAMES | google.com | Host name to query for; seperate by comma for multiple |
| APP_LOGGLY_TOKEN | - | Optional; Loggly token to log to Loggly |

## Metrics

Metrics are exposed on port 8080, the following dns specific metric is exposed

`dc_dns_request_total{app=[application_name], found=[true or false], nameServer=[name server address], targetHost=[target host address]} [number of requests]`

Example:

`dc_dns_requests_total{app="DnsChecker-Local",found="true",nameServer="1.1.1.1",targetHost="google.com"} 3`