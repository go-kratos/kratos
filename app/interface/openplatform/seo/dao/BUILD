package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "item.go",
        "pro.go",
        "sitemap.go",
    ],
    importpath = "go-common/app/interface/openplatform/seo/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/openplatform/seo/conf:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/log:go_default_library",
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
