.PHONY: build clean test install uninstall

BIN_DIR := bin
PREFIX := /usr/local
INSTALL_DIR := $(PREFIX)/bin
BINARIES := templater

build: $(addprefix $(BIN_DIR)/,$(BINARIES))

$(BIN_DIR)/templater: main.go $(shell find internal -name '*.go')
	@mkdir -p $(BIN_DIR)
	go build -o $@ .

clean:
	rm -rf $(BIN_DIR)

test:
	go test ./... -count=1

install: build
	@mkdir -p $(INSTALL_DIR)
	@for bin in $(BINARIES); do \
		cp $(BIN_DIR)/$$bin $(INSTALL_DIR)/$$bin; \
		chmod 755 $(INSTALL_DIR)/$$bin; \
	done
	@echo "Installed to $(INSTALL_DIR)"

uninstall:
	@for bin in $(BINARIES); do \
		rm -f $(INSTALL_DIR)/$$bin; \
	done
	@echo "Uninstalled from $(INSTALL_DIR)"