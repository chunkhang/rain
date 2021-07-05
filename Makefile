all: rain

rain: main.go
	go build -o $@

.PHONY: all
