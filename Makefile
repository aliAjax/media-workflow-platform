.PHONY: test lint run lines smoke
test:
	go test ./...
lint:
	gofmt -w $$(find . -name '*.go')
	go vet ./...
run:
	go run ./cmd/media-api
lines:
	./tools/count-lines.sh
smoke:
	./scripts/smoke.sh
