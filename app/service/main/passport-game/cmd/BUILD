load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    data = ["passport-game-service.toml"],
    importpath = "go-common/app/service/main/passport-game/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
    deps = [
        "//app/service/main/passport-game/conf:go_default_library",
        "//app/service/main/passport-game/http:go_default_library",
        "//app/service/main/passport-game/service:go_default_library",
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
