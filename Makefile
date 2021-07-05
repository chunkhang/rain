all: rain

rain: *.go
	go build -o $@

.PHONY: all
