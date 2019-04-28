package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "apk.go",
        "helper.go",
        "unzip.go",
        "writer.go",
    ],
    importpath = "go-common/app/admin/main/macross/tools",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//library/log:go_default_library"],
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
