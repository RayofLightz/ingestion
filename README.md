# Ingestion

Ingestion is a a multithreaded stream processor built to enrich suricata logs

## Installation

### Basic Installation

Use go

```bash
go get -d github.com/RayofLightz/ingestion
```
Then build ingestion like so
```bash
cd $GOPATH/src/github.com/RayofLightz/ingestion;
go build
```

### With malware host name detection
To install the malware host detection csv data set run `scripts/downloadhosts.sh` and then edit the config file.


## Usage

```bash
./ingestion
```

## Configuration

Ingestion can be configured using the config file `config/config.json`.
Currently there are three configuation values.

### local
The local configuration option has ingestion run only on local host when set to true. When set to false it binds to 0.0.0.0

### rev_lookup
When this configuration value is set to true ingestion dose a reverse lookup for the `dest_ip`.

### check_known_malware
When set to true uses the result from rev_lookup against a set of known malicious domains. REV_LOOKUP MUST BE SET TO TRUE FOR THIS TO BE SET TO TRUE.

## Sending ingestion data
Currently ingestion has data sent one log entry per tcp request. Suricata seperates json entries using a new line the current way to send a json entry is `head -1 eve.json | nc 127.0.0.1 8080`. This will send the the first json entry to ingestion.
The next thing I plan on implimenting is a tool to send ingestion log data.

## Contributing
Pull requests are welcome. I am not picky about style, but please run the script included in the scripts dir like so `scripts/massfmt.sh` before commiting. 

## License
[BSD](https://choosealicense.com/licenses/bsd/)
