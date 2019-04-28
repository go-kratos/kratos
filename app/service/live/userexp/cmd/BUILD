load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    data = ["live-userexp-service.toml"],
    importpath = "go-common/app/service/live/userexp/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
    deps = [
        "//app/service/live/userexp/conf:go_default_library",
        "//app/service/live/userexp/http:go_default_library",
        "//app/service/live/userexp/service:go_default_library",
        "//library/log:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/syscall:go_default_library",
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
