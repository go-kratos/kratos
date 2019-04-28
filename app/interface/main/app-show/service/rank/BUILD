package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["rank_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-show/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "audit.go",
        "rank.go",
    ],
    importpath = "go-common/app/interface/main/app-show/service/rank",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-show/conf:go_default_library",
        "//app/interface/main/app-show/dao/account:go_default_library",
        "//app/interface/main/app-show/dao/archive:go_default_library",
        "//app/interface/main/app-show/dao/audit:go_default_library",
        "//app/interface/main/app-show/dao/recommend:go_default_library",
        "//app/interface/main/app-show/dao/region:go_default_library",
        "//app/interface/main/app-show/dao/relation:go_default_library",
        "//app/interface/main/app-show/model:go_default_library",
        "//app/interface/main/app-show/model/region:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/relation/model:go_default_library",
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
