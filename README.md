# vigor-exporter

A (crusty) DrayTek Vigor 130 (and others) exporter for prometheus.

This can be used to scrape your Vigor modem for some connection metrics.

## Build

```
go build .
```

## Usage Information

### Run directly

```
./vigor-exporter --help
  -host string
    	hostname/ip the Vigor is reachable on
  -password string
    	password to authenticate to the Vigor
  -username string
    	username to authenticate to the Vigor
```

The exporter is listening on `*:9103` and provides metrics at the `/metrics`
path.

### Run with docker

```
docker build -t vigor  .

docker run vigor
  -host string
    	hostname/ip the Vigor is reachable on
  -password string
    	password to authenticate to the Vigor
  -username string
    	username to authenticate to the Vigor
```
