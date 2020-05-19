
.PHONY: all build exec clean

all: build exec clean


docker:
	@GOOS=linux GOARCH=amd64 ./build
	docker build -t raylax/imx:0.0.1 .
	@rm ./bin/imx

build:
	@GO_BUILD_FLAGS="-v" ./build

exec:
	@./bin/imx -v

clean:
	@rm ./bin/imx