load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    data = ["ticket-sales.toml"],
    importpath = "go-common/app/service/openplatform/ticket-sales/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
    deps = [
        "//app/service/openplatform/ticket-sales/conf:go_default_library",
        "//app/service/openplatform/ticket-sales/server/grpc:go_default_library",
        "//app/service/openplatform/ticket-sales/server/http:go_default_library",
        "//app/service/openplatform/ticket-sales/service:go_default_library",
        "//app/service/openplatform/ticket-sales/service/mis:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/ecode/tip:go_default_library",
        "//library/log:go_default_library",
        "//library/net/trace:go_default_library",
    ],
)

go_binary(
    name = "cmd",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
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
