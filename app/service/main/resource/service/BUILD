package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "ads_test.go",
        "banner_test.go",
        "relate_test.go",
        "resource_test.go",
        "service_test.go",
        "special_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/resource/api/v1:go_default_library",
        "//app/service/main/resource/conf:go_default_library",
        "//library/cache/redis:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "abtest.go",
        "ads.go",
        "archive.go",
        "banner.go",
        "relate.go",
        "resource.go",
        "service.go",
        "show.go",
        "special.go",
    ],
    importpath = "go-common/app/service/main/resource/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/location/model:go_default_library",
        "//app/service/main/location/rpc/client:go_default_library",
        "//app/service/main/resource/api/v1:go_default_library",
        "//app/service/main/resource/conf:go_default_library",
        "//app/service/main/resource/dao/abtest:go_default_library",
        "//app/service/main/resource/dao/ads:go_default_library",
        "//app/service/main/resource/dao/alarm:go_default_library",
        "//app/service/main/resource/dao/cpm:go_default_library",
        "//app/service/main/resource/dao/manager:go_default_library",
        "//app/service/main/resource/dao/resource:go_default_library",
        "//app/service/main/resource/dao/show:go_default_library",
        "//app/service/main/resource/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/xstr:go_default_library",
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
