# DnsChecker

Small application to query specific dns servers for host names.

## Usage

In order to run DnsChecker you have three options;

### Docker

``` bash
docker run mkromkamp/dnschecker:latest
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
