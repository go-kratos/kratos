package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "resolver.go",
        "util.go",
    ],
    importpath = "go-common/library/net/rpc/warden/resolver",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/conf/env:go_default_library",
        "//library/log:go_default_library",
        "//library/naming:go_default_library",
        "//library/net/rpc/warden/metadata:go_default_library",
        "//vendor/github.com/dgryski/go-farm:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "@org_golang_google_grpc//resolver:go_default_library",
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
        "//library/net/rpc/warden/resolver/direct:all-srcs",
        "//library/net/rpc/warden/resolver/livezk:all-srcs",
        "//library/net/rpc/warden/resolver/test:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
