# Copyright 2022 shun159 <dreamdiagnosis@gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

CMD_GO ?= go

.PHONY: all

all: get-deps test-unit

get-deps:
	$(CMD_GO) mod download
	$(CMD_GO) mod tidy

GO_UTEST_FILES:=$(shell find . -type f -name '*_test.go' -print)
.PHONY: test-unit
test-unit:
	$(CMD_GO) test $(GO_UTEST_FILES) -v