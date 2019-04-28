package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["tag_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-tag/conf:go_default_library",
        "//app/interface/main/app-tag/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "cache.go",
        "infoc.go",
        "service.go",
        "tag.go",
    ],
    importpath = "go-common/app/interface/main/app-tag/service/tag",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-tag/conf:go_default_library",
        "//app/interface/main/app-tag/dao/archive:go_default_library",
        "//app/interface/main/app-tag/dao/bangumi:go_default_library",
        "//app/interface/main/app-tag/dao/recommend:go_default_library",
        "//app/interface/main/app-tag/dao/region:go_default_library",
        "//app/interface/main/app-tag/dao/tag:go_default_library",
        "//app/interface/main/app-tag/model:go_default_library",
        "//app/interface/main/app-tag/model/bangumi:go_default_library",
        "//app/interface/main/app-tag/model/feed:go_default_library",
        "//app/interface/main/app-tag/model/region:go_default_library",
        "//app/interface/main/app-tag/model/tag:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/stat/prom:go_default_library",
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
