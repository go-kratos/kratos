package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "databus.go",
        "load.go",
        "search.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/riot-search/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/riot-search/conf:go_default_library",
        "//app/service/main/riot-search/dao:go_default_library",
        "//app/service/main/riot-search/model:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//vendor/github.com/ivpusic/grpool:go_default_library",
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
