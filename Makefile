t:
	go test ./...
b:
	TEST_MODE="false" go run .
d:
	TEST_MODE="false" dlv debug -l 127.0.0.1:9080 --headless --accept-multiclient --api-version 2 --output ./debug --continue
