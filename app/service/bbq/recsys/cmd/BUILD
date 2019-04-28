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
    data = ["test.toml"],
    importpath = "go-common/app/service/bbq/recsys/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/recsys/conf:go_default_library",
        "//app/service/bbq/recsys/server/grpc:go_default_library",
        "//app/service/bbq/recsys/server/http:go_default_library",
        "//app/service/bbq/recsys/service:go_default_library",
        "//library/ecode/tip:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/service/bbq/recsys/cmd/client:all-srcs",
        "//app/service/bbq/recsys/cmd/relatedrecsys:all-srcs",
        "//app/service/bbq/recsys/cmd/upsrecsys:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
