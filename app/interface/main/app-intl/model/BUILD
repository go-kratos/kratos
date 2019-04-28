package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["const.go"],
    importpath = "go-common/app/interface/main/app-intl/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
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
        "//app/interface/main/app-intl/model/bangumi:all-srcs",
        "//app/interface/main/app-intl/model/feed:all-srcs",
        "//app/interface/main/app-intl/model/manager:all-srcs",
        "//app/interface/main/app-intl/model/player:all-srcs",
        "//app/interface/main/app-intl/model/region:all-srcs",
        "//app/interface/main/app-intl/model/search:all-srcs",
        "//app/interface/main/app-intl/model/tag:all-srcs",
        "//app/interface/main/app-intl/model/view:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
