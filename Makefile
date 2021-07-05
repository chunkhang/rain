all: rain

rain: main.go
	go build -ldflags -a -o $@

.PHONY: all
