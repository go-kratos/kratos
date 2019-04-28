package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["wrr_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/net/metadata:go_default_library",
        "//library/net/rpc/warden/metadata:go_default_library",
        "@org_golang_google_grpc//balancer:go_default_library",
        "@org_golang_google_grpc//resolver:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["wrr.go"],
    importpath = "go-common/library/net/rpc/warden/balancer/wrr",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/rpc/warden/metadata:go_default_library",
        "//library/stat/summary:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//balancer:go_default_library",
        "@org_golang_google_grpc//balancer/base:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_grpc//resolver:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
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
        "//library/net/rpc/warden/balancer/wrr/test:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
