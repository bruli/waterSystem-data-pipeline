module github.com/bruli/waterSystem-data-pipeline

go 1.26.1

require (
	github.com/caarlos0/env/v11 v11.4.0
	github.com/nats-io/nats.go v1.49.0
)

require (
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/nats-io/nkeys v0.4.12 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
)

tool (
	github.com/mfridman/tparse
	golang.org/x/vuln/cmd/govulncheck
	mvdan.cc/gofumpt
)