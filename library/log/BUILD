package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "agent_test.go",
        "dsn_test.go",
        "encode_test.go",
        "log_test.go",
        "pattern_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/log/internal:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "agent.go",
        "doc.go",
        "dsn.go",
        "file.go",
        "handler.go",
        "level.go",
        "log.go",
        "logrus.go",
        "pattern.go",
        "stdout.go",
        "util.go",
        "verbose.go",
    ],
    importpath = "go-common/library/log",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/conf/dsn:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/log/internal:go_default_library",
        "//library/log/internal/filewriter:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/sirupsen/logrus:go_default_library",
    ],
)

go_test(
    name = "go_default_xtest",
    srcs = ["example_test.go"],
    tags = ["automanaged"],
    deps = ["//library/log:go_default_library"],
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
        "//library/log/anticheat:all-srcs",
        "//library/log/benchmark:all-srcs",
        "//library/log/infoc:all-srcs",
        "//library/log/internal:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
