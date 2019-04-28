package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "jury_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/credit/conf:go_default_library",
        "//app/admin/main/credit/model/blocked:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "blocked.go",
        "extra_func.go",
        "jury.go",
        "labour.go",
        "msg.go",
        "service.go",
        "user.go",
    ],
    importpath = "go-common/app/admin/main/credit/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/credit/conf:go_default_library",
        "//app/admin/main/credit/dao/account:go_default_library",
        "//app/admin/main/credit/dao/blocked:go_default_library",
        "//app/admin/main/credit/dao/manager:go_default_library",
        "//app/admin/main/credit/dao/msg:go_default_library",
        "//app/admin/main/credit/dao/relation:go_default_library",
        "//app/admin/main/credit/dao/search:go_default_library",
        "//app/admin/main/credit/dao/upload:go_default_library",
        "//app/admin/main/credit/model:go_default_library",
        "//app/admin/main/credit/model/blocked:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/coin/api:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
