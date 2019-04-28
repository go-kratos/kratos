# Copyright (C) 2015  The GoHBase Authors.  All rights reserved.
# This file is part of GoHBase.
# Use of this source code is governed by the Apache License 2.0
# that can be found in the COPYING file.

GO := go
TEST_TIMEOUT := 30s
INTEGRATION_TIMEOUT := 120s
GOTEST_FLAGS := -v

GOLINT := golint

all: install

install:
	$(GO) install ./...

check: vet fmtcheck lint test
jenkins: check integration

COVER_PKGS := `find ./* -name '*_test.go' ! -path "./test/mock/*" | xargs -I{} dirname {} | sort -u`
COVER_MODE := count
coverdata:
	echo 'mode: $(COVER_MODE)' >coverage.out
	for dir in $(COVER_PKGS); do \
	  $(GO) test -covermode=$(COVER_MODE) -tags=testing -coverprofile=cov.out-t $$dir || exit; \
	  tail -n +2 cov.out-t >> coverage.out && \
	  rm cov.out-t; \
	done;

coverage: coverdata
	$(GO) tool cover -html=coverage.out
	rm -f coverage.out

fmtcheck:
	errors=`gofmt -l .`; if test -n "$$errors"; then echo Check these files for style errors:; echo "$$errors"; exit 1; fi
	find . -name '*.go' ! -path "./pb/*" ! -path "./test/mock/*" !  -path './gen.go' -exec ./check_line_len.awk {} +

vet:
	$(GO) vet ./...

lint:
	find ./* -type d ! -name pb ! -name mock ! -path "./test/mock/*" | xargs -L 1 $(GOLINT) &>lint; :
	if test -s lint; then echo Check these packages for golint:; cat lint; rm lint; exit 1; else rm lint; fi
# The above is ugly, but unfortunately golint doesn't exit 1 when it finds
# lint.  See https://github.com/golang/lint/issues/65

test:
	$(GO) test $(GOTEST_FLAGS) -race -timeout=$(TEST_TIMEOUT) -tags=testing ./...

integration:
	$(GO) test $(GOTEST_FLAGS) -race -timeout=$(INTEGRATION_TIMEOUT) -tags=integration

.PHONY: all check coverage coverdata fmtcheck install integration jenkins lint test vet
