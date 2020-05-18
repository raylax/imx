
.PHONY: all build exec clean

all: build exec clean

build:
	@GO_BUILD_FLAGS="-v" ./build

exec:
	@./bin/imx -v

clean:
	@rm ./bin/imx