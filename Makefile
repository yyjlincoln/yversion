.PHONY: all
all: build

.PHONY: build
build:
	@go build

.PHONY: clean
clean:
	@rm ./yversion
