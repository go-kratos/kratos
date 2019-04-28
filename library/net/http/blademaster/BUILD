package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "client_test.go",
        "server_test.go",
        "trace_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/net/http/blademaster/tests:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/netutil/breaker:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "context.go",
        "cors.go",
        "csrf.go",
        "device.go",
        "logger.go",
        "metadata.go",
        "perf.go",
        "prometheus.go",
        "recovery.go",
        "routergroup.go",
        "server.go",
        "trace.go",
        "utils.go",
    ],
    importpath = "go-common/library/net/http/blademaster",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/conf/dsn:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/net/ip:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/netutil/breaker:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/stat:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/prometheus/client_golang/prometheus/promhttp:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_gogo_protobuf//types:go_default_library",
    ],
)

go_test(
    name = "go_default_xtest",
    srcs = ["example_test.go"],
    tags = ["automanaged"],
    deps = [
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/http/blademaster/tests:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//library/net/http/blademaster/binding:all-srcs",
        "//library/net/http/blademaster/middleware/antispam:all-srcs",
        "//library/net/http/blademaster/middleware/auth:all-srcs",
        "//library/net/http/blademaster/middleware/cache:all-srcs",
        "//library/net/http/blademaster/middleware/limit/aqm:all-srcs",
        "//library/net/http/blademaster/middleware/permit:all-srcs",
        "//library/net/http/blademaster/middleware/proxy:all-srcs",
        "//library/net/http/blademaster/middleware/rate:all-srcs",
        "//library/net/http/blademaster/middleware/supervisor:all-srcs",
        "//library/net/http/blademaster/middleware/tag:all-srcs",
        "//library/net/http/blademaster/middleware/verify:all-srcs",
        "//library/net/http/blademaster/render:all-srcs",
        "//library/net/http/blademaster/tests:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
