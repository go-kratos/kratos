package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["server_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/location/api:go_default_library",
        "//app/service/main/location/conf:go_default_library",
        "//app/service/main/location/service:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/time:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["server.go"],
    importpath = "go-common/app/service/main/location/server/grpc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/location/api:go_default_library",
        "//app/service/main/location/model:go_default_library",
        "//app/service/main/location/service:go_default_library",
        "//library/net/rpc/warden:go_default_library",
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
