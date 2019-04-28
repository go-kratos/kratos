package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "region_test.go",
        "show_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-show/conf:go_default_library",
        "//app/interface/main/app-show/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "audit.go",
        "banner.go",
        "cache.go",
        "infoc.go",
        "region.go",
        "service.go",
        "show.go",
    ],
    importpath = "go-common/app/interface/main/app-show/service/region",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-show/conf:go_default_library",
        "//app/interface/main/app-show/dao/activity:go_default_library",
        "//app/interface/main/app-show/dao/ad:go_default_library",
        "//app/interface/main/app-show/dao/archive:go_default_library",
        "//app/interface/main/app-show/dao/audit:go_default_library",
        "//app/interface/main/app-show/dao/bangumi:go_default_library",
        "//app/interface/main/app-show/dao/card:go_default_library",
        "//app/interface/main/app-show/dao/dynamic:go_default_library",
        "//app/interface/main/app-show/dao/location:go_default_library",
        "//app/interface/main/app-show/dao/recommend:go_default_library",
        "//app/interface/main/app-show/dao/region:go_default_library",
        "//app/interface/main/app-show/dao/resource:go_default_library",
        "//app/interface/main/app-show/dao/search:go_default_library",
        "//app/interface/main/app-show/dao/tag:go_default_library",
        "//app/interface/main/app-show/model:go_default_library",
        "//app/interface/main/app-show/model/bangumi:go_default_library",
        "//app/interface/main/app-show/model/banner:go_default_library",
        "//app/interface/main/app-show/model/card:go_default_library",
        "//app/interface/main/app-show/model/region:go_default_library",
        "//app/interface/main/app-show/model/tag:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/location/model:go_default_library",
        "//app/service/main/resource/model:go_default_library",
        "//app/service/openplatform/pgc-season/api/grpc/season/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/dgryski/go-farm:go_default_library",
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
