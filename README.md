# InBrief scraper

InBrief scraper uses TDLib and gRPC to access messages and allow to perform
scraping messages offline (using RPC/HTTP) or online (by subscribing to Redis)
events

## Requirements

- Go 1.24.3 or higher
- Redis 9.0 or higher
- TDLib (compiled with Golang bindings)
- AWS CLI (optional, for AWS services)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/nrydanov/inbrief.git
cd inbrief
```

2. Install dependencies:
```bash
go mod download
```

3. Generate proto files:
```bash
buf generate
```

## Running

To run in development mode:
```bash
air
```

## License

MIT
