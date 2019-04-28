load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
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
    srcs = [
        "archive.go",
        "copyright.go",
        "es.go",
        "log.go",
        "mapdata.go",
        "mng.go",
        "mng_v2.go",
        "param.go",
        "query.go",
        "query_extra.go",
        "uname.go",
        "upsert.go",
    ],
    importpath = "go-common/app/admin/main/search/model",
    tags = ["automanaged"],
    deps = [
        "//library/log:go_default_library",
        "//vendor/gopkg.in/olivere/elastic.v5:go_default_library",
    ],
)
