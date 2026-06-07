APP := pgpackage
CMD := ./cmd/pgpackage
SAMPLE_PROJECT := testdata/sample/sample.pgpackage
OUT_DIR := out
SAMPLE_PACKAGE := $(OUT_DIR)/SampleProject.pgpkg

.PHONY: build test sample package clean

build:
	go build -o $(APP) $(CMD)

test:
	go test ./...

sample: $(SAMPLE_PACKAGE)

package: $(SAMPLE_PACKAGE)

$(SAMPLE_PACKAGE):
	go run $(CMD) build --project $(SAMPLE_PROJECT) --output $(OUT_DIR)/

clean:
	rm -rf $(OUT_DIR) $(APP)
