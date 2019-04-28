package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["view.go"],
    importpath = "go-common/app/interface/main/app-intl/model/view",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-intl/model:go_default_library",
        "//app/interface/main/app-intl/model/bangumi:go_default_library",
        "//app/interface/main/app-intl/model/manager:go_default_library",
        "//app/interface/main/app-intl/model/tag:go_default_library",
        "//app/interface/main/dm2/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
