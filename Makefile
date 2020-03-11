
fmt:
	find . -name '*.go' -exec gofmt -s -w {} \;
	go build ./config/ini
	go build