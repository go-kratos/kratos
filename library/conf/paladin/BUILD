package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "default.go",
        "file.go",
        "helper.go",
        "map.go",
        "mock.go",
        "sven.go",
        "toml.go",
        "value.go",
    ],
    importpath = "go-common/library/conf/paladin",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/ip:go_default_library",
        "//library/net/netutil:go_default_library",
        "//vendor/github.com/BurntSushi/toml:go_default_library",
        "//vendor/github.com/naoina/toml:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
    ],
)

go_test(
    name = "go_default_xtest",
    srcs = [
        "example_test.go",
        "file_test.go",
        "map_test.go",
        "mock_test.go",
    ],
    tags = ["automanaged"],
    deps = [
        "//library/conf/paladin:go_default_library",
        "//vendor/github.com/naoina/toml:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "sven_test.go",
        "value_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//library/conf/env:go_default_library",
        "//vendor/github.com/naoina/toml:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
    ],
)
