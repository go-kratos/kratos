load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

proto_library(
    name = "liverpc_proto",
    srcs = ["liverpc.proto"],
    tags = ["automanaged"],
)

go_proto_library(
    name = "liverpc_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "go-common/library/net/rpc/liverpc",
    proto = ":liverpc_proto",
    tags = ["automanaged"],
)

go_test(
    name = "go_default_test",
    srcs = ["client_conn_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/conf/env:go_default_library",
        "//library/net/metadata:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "client_conn.go",
        "option.go",
        "protocol.go",
    ],
    embed = [":liverpc_go_proto"],
    importpath = "go-common/library/net/rpc/liverpc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/conf/env:go_default_library",
        "//library/log:go_default_library",
        "//library/naming:go_default_library",
        "//library/naming/discovery:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/stat:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/json-iterator/go:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
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
        "//library/net/rpc/liverpc/context:all-srcs",
        "//library/net/rpc/liverpc/test:all-srcs",
        "//library/net/rpc/liverpc/testdata:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
