load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)
load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

package(default_visibility = ["//visibility:public"])

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
        "//library/queue/databus/databusutil:all-srcs",
        "//library/queue/databus/metadata:all-srcs",
        "//library/queue/databus/report:all-srcs",
    ],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["databus.go"],
    embed = [":databus_go_proto"],
    importpath = "go-common/library/queue/databus",
    tags = ["automanaged"],
    deps = [
        "//library/cache/redis:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/container/pool:go_default_library",
        "//library/log:go_default_library",
        "//library/naming:go_default_library",
        "//library/naming/discovery:go_default_library",
        "//library/net/netutil:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
    ],
)

go_test(
    name = "go_default_xtest",
    srcs = ["databus_test.go"],
    tags = ["automanaged"],
    deps = [
        "//library/naming/discovery:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/time:go_default_library",
    ],
)

proto_library(
    name = "databus_proto",
    srcs = ["databus.proto"],
    tags = ["automanaged"],
    deps = ["@gogo_special_proto//github.com/gogo/protobuf/gogoproto"],
)

go_proto_library(
    name = "databus_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_proto"],
    importpath = "go-common/library/queue/databus",
    proto = ":databus_proto",
    tags = ["automanaged"],
    deps = ["@com_github_gogo_protobuf//gogoproto:go_default_library"],
)
