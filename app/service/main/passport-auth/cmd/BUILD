package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "cmd",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    data = [
        "passport-auth-service.toml",
        "passport-auth-service-uat.toml",
    ],
    importpath = "go-common/app/service/main/passport-auth/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/passport-auth/conf:go_default_library",
        "//app/service/main/passport-auth/http:go_default_library",
        "//app/service/main/passport-auth/rpc/grpc:go_default_library",
        "//app/service/main/passport-auth/rpc/server:go_default_library",
        "//app/service/main/passport-auth/service:go_default_library",
        "//library/log:go_default_library",
        "//library/net/trace:go_default_library",
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
