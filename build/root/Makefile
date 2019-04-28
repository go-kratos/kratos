# Don't allow an implicit 'all' rule. This is not a user-facing file.
ifeq ($(MAKECMDGOALS),)
    $(error This Makefile requires an explicit rule to be specified)
endif

ifeq ($(DBG_MAKEFILE),1)
    $(warning ***** starting Makefile.generated_files for goal(s) "$(MAKECMDGOALS)")
    $(warning ***** $(shell date))
endif
GOBIN  := $(go env GOBIN)
ifeq ($(GOBIN),)
   GOBIN := ~/go/bin
endif
# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /bin/bash
ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"
# We don't need make's built-in rules.
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

# Constants used throughout.
.EXPORT_ALL_VARIABLES:
OUT_DIR ?= _output
BIN_DIR := $(OUT_DIR)/bin

.PHONY: build update clean
all: check bazel-update 

build: init check bazel-build 

test-coverage:
	bazel coverage  --test_env=DEPLOY_ENV=uat --test_timeout=60 --test_env=APP_ID=bazel.test --test_output=all --cache_test_results=no   //app/service/main/account/dao/...

simple-build:
	bazel build --watchfs -- //tools/... -//vendor/...
ifeq ($(WHAT),)
bazel-build:
	bazel build --config=office --watchfs //app/... //build/... //library/... 
else 
bazel-build: 
	bazel build --config=ci -- //$(WHAT)/...
endif

build-keep-going:
	bazel build --config=ci -k //app/... //build/... //library/... 
	cat bazel-out/stable-status.txt
clean:
	bazel clean --expunge 
	rm -rf _output
update: init bazel-update

bazel-update:
	./build/update-bazel.sh
prow-update:
	./build/update-prow.sh
test:
	@if [ "$(WHAT)" !=  "" ]; \
         then \
	 cd $(WHAT) && make ; \
	 else \
	 echo "Please input the WHAT" ;\
	 fi

bazel-test:
	@if [ "$(WHAT)" !=  "" ]; \
         then \
	 bazel test --watchfs -- //$(WHAT)/... ; \
	 else \
	 echo "Please input the WHAT" ;\
	 fi
check:
	@./build/check.sh
init:
	@if [ ! -f .git/hooks/pre-commit ] ; \
	then \
	echo "make all" >> .git/hooks/pre-commit; \
	sudo chmod +x .git/hooks/pre-commit; \
	fi
build-all-kratos:
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_386  //app/tool/kratos:kratos
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //app/tool/kratos:kratos
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:darwin_amd64 //app/tool/kratos:kratos


install-kratos: init check build-kratos
	@if [[ "$(ARCH)" == "Linux" ]]; then \
		cp bazel-bin/app/tool/kratos/linux_amd64_pure_stripped/kratos $(GOBIN); \
	fi; \
	if [[ "$(ARCH)" == "Darwin" ]]; then \
		cp bazel-bin/app/tool/kratos/darwin_amd64_stripped/kratos $(GOBIN); \
	fi

build-kratos:
	bazel build  //app/tool/kratos:kratos 

ci-bazel-build: 
	bazel build --config=ci -- //app/...

ci-bazel-build-a: 
	bazel build --config=ci -- //app/admin/...

ci-bazel-build-b: 
	bazel build --config=ci -- //app/interface/...

ci-bazel-build-c: 
	bazel build --config=ci -- //app/job/... //app/tool/... //app/common/... //app/infra/...

ci-bazel-build-d: 
	bazel build --config=ci -- //app/service/... //library/...

ci-bazel-build-common: 
	bazel build --config=ci -- //app/common/...

ci-bazel-build-infra: 
	bazel build --config=ci -- //app/infra/...

ci-bazel-build-tool: 
	bazel build --config=ci -- //app/tool/...

ci-bazel-build-main:
	bazel build --config=ci -- //app/admin/main/... //app/interface/main/... //app/job/main/... //app/service/main/... 

ci-bazel-build-live:
	bazel build --config=ci -- //app/admin/live/... //app/interface/live/... //app/job/live/... //app/job/live-userexp/... //app/service/live/... 

ci-bazel-build-ep:
	bazel build --config=ci -- //app/admin/ep/... //app/service/ep/... 

ci-bazel-build-openplatform:
	bazel build --config=ci -- //app/admin/openplatform/... //app/interface/openplatform/... //app/job/openplatform/... //app/service/openplatform/... 

ci-bazel-build-bbq:
	bazel build --config=ci -- //app/interface/bbq/... //app/job/bbq/... //app/service/bbq/... 

ci-bazel-build-video:
	bazel build --config=ci -- //app/interface/video/... //app/service/video/... 

ci-bazel-build-ops:
	bazel build --config=ci -- //app/service/ops/...

ci-bazel-build-library: 
	bazel build --config=ci -- //library/...

ci-bazel-build-admin-main:
	bazel build --config=ci -- //app/admin/main/... 

ci-bazel-build-interface-main:
	bazel build --config=ci -- //app/interface/main/... 

ci-bazel-build-job-main:
	bazel build --config=ci -- //app/job/main/... 

ci-bazel-build-service-main:
	bazel build --config=ci -- //app/service/main/... 
