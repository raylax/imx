
.PHONY: build
build:
	@GO_BUILD_FLAGS="-v" ./build
	@./bin/imx -v

clean:
	@rm ./bin/imx