package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["search_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-intl/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["search.go"],
    importpath = "go-common/app/interface/main/app-intl/service/search",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-interface/model/tag:go_default_library",
        "//app/interface/main/app-intl/conf:go_default_library",
        "//app/interface/main/app-intl/dao/account:go_default_library",
        "//app/interface/main/app-intl/dao/archive:go_default_library",
        "//app/interface/main/app-intl/dao/article:go_default_library",
        "//app/interface/main/app-intl/dao/bangumi:go_default_library",
        "//app/interface/main/app-intl/dao/resource:go_default_library",
        "//app/interface/main/app-intl/dao/search:go_default_library",
        "//app/interface/main/app-intl/dao/tag:go_default_library",
        "//app/interface/main/app-intl/model:go_default_library",
        "//app/interface/main/app-intl/model/bangumi:go_default_library",
        "//app/interface/main/app-intl/model/search:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
