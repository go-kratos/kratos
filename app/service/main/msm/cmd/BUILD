load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    data = ["msm-service-example.toml"],
    importpath = "go-common/app/service/main/msm/cmd",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/msm/conf:go_default_library",
        "//app/service/main/msm/http:go_default_library",
        "//app/service/main/msm/service:go_default_library",
        "//library/log:go_default_library",
    ],
)

go_binary(
    name = "cmd",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)
