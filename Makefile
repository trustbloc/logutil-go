# Copyright SecureKey Technologies Inc.
#
# SPDX-License-Identifier: Apache-2.0

GOBIN_PATH=$(abspath .)/.build/bin

# Tool commands (overridable)
GO_VER ?= 1.25

OS := $(shell uname)
ifeq  ($(OS),$(filter $(OS),Darwin Linux))
	PATH:=$(PATH):$(GOBIN_PATH)
else
	PATH:=$(PATH);$(subst /,\\,$(GOBIN_PATH))
endif

export GOTOOLCHAIN=go1.25.3+auto

.PHONY: all
all: clean checks unit-test

.PHONY: checks
checks: license #lint

.PHONY: lint
lint:
	@scripts/check_lint.sh

.PHONY: license
license:
	@scripts/check_license.sh

.PHONY: unit-test
unit-test:
	@scripts/check_unit.sh

.PHONY: clean
clean:
	@rm -rf ./build
	@rm -rf coverage*.out
