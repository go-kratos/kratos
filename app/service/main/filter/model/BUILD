package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "ai.go",
        "area.go",
        "filter.go",
        "hit.go",
        "http.go",
        "key.go",
        "limit.go",
        "model.go",
        "white.go",
    ],
    importpath = "go-common/app/service/main/filter/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/filter/conf:go_default_library",
        "//app/service/main/filter/model/actriearea:go_default_library",
        "//app/service/main/filter/model/actriekey:go_default_library",
        "//app/service/main/filter/service/regexp:go_default_library",
        "//library/time:go_default_library",
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
        "//app/service/main/filter/model/actriearea:all-srcs",
        "//app/service/main/filter/model/actriekey:all-srcs",
        "//app/service/main/filter/model/lrulist:all-srcs",
        "//app/service/main/filter/model/rpc:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
