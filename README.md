# super-duper-succotash

---

## Building Instruction

### Requirements

- go 1.13

Run following commands after clone this repository and cd into it

```sh
go build ./cmd/ -o main
```

## Configurations and Executing

This excutable has 2 mode: `oneshot` and `service` which can be selected via `-m` or `--mode`
please see help message below, or run with `--help` flag to print the following help message

```
Usage:
  main [OPTIONS]

Application Options:
  -a, --amount=                input amount to calculate
  -i, --input-asset=           input asset type, output asset type will be automatically set via pair config according to exchange engine, if available
  -o, --output-asset=          output asset type, can be set if engine support exchange routing with more than 1 pair
  -m, --mode=[oneshot|service] select wheter to run as oneshot or until manually stop
  -E, --engine=[coinbase_pro]  select exchange engine to use
  -e, --engine-config=         configuration for exchange engine, in key:value format, one pair per each flag

Help Options:
  -h, --help                   Show this help message
```

#### Engine Configurations

Currently only `coinbase_pro` can be used as engine, which also has its own configurations
Supply each key:value with repeating `-e` or `--engine-config` flag.
Output asset is automatically set by trading engine, flag is ignored
Input asset value is case-insensitive
example below.

```sh
# run as service
./main -m service -E 'coinbase_pro' -e 'api_url:https://api.pro.coinbase.com' -e 'api_level:1' -e 'pair:ETH-USD' -e 'poll_interval:5s' -a "1000000" -i "eth"
# or run only once
./main -m oneshot -E 'coinbase_pro' -e 'api_url:https://api.pro.coinbase.com' -e 'api_level:2' -e 'pair:BTC-USD' -e 'poll_interval:5s' -a "1000000" -i "btc"
```

###### Configuration type signature for `coinbase_pro` engine

```go
// Engine holds necessary configuration for coinbase pro API
type Engine struct {
	APIURL       string        `mapstructure:"api_url"`
	APILevel     int64         `mapstructure:"api_level"`
	Pair         string        `mapstructure:"pair"`
	PollInterval time.Duration `mapstructure:"poll_interval"`
}
```

## Design Rationale

The API will not attempt to handle **transient error**, when found one, it will most likely be panicked
this enabled us to restart promptly after service is crashed, service supervisor such as kubernetes, or init system such as systemd is recommended

Using weak interface allow us to have custom configuration for each engine that will be implemented in the future,
the trade-off is that it will be harder to debug such an architecture when something went wrong, so error handling with stacktrace is must.

There is various improvement that can be made, such as:

- custom error type and error propagations
- instrumenting, logging and tracing to external service
- pseudo enum type, which is quite cubersome to implement in go since it has no proper sum type
- integration test and test for function that can be panicked
- 12 factor app configuration model
- output format other than stdout for human, e.g. message queue, JSON stream and websocket.

This project is implemented with incremental refactor-able in mind, and sacrifice some part of simplicity for adaptability
because the author may use this for other purpose in the future, and effort that has made should not be wasted.
nonetheless, this is still just a minimum viable product
