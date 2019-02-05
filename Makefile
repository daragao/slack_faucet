LD_FLAGS := -ldflags "-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git rev-parse HEAD`"

run:
	go run $(LD_FLAGS) main.go

build:
	go build -o faucet $(LD_FLAGS) main.go
