package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["service_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/creative/conf:go_default_library",
        "//app/interface/main/creative/model/archive:go_default_library",
        "//app/interface/main/creative/model/order:go_default_library",
        "//app/interface/main/creative/service:go_default_library",
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "biz.go",
        "fav.go",
        "history.go",
        "mission.go",
        "myinfo.go",
        "netsafe.go",
        "order.go",
        "porder.go",
        "rules.go",
        "service.go",
        "staff.go",
        "tag.go",
        "viewpoint.go",
    ],
    importpath = "go-common/app/interface/main/creative/service/archive",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/videoup/model/archive:go_default_library",
        "//app/interface/main/creative/conf:go_default_library",
        "//app/interface/main/creative/dao/account:go_default_library",
        "//app/interface/main/creative/dao/activity:go_default_library",
        "//app/interface/main/creative/dao/appeal:go_default_library",
        "//app/interface/main/creative/dao/archive:go_default_library",
        "//app/interface/main/creative/dao/coin:go_default_library",
        "//app/interface/main/creative/dao/game:go_default_library",
        "//app/interface/main/creative/dao/order:go_default_library",
        "//app/interface/main/creative/dao/search:go_default_library",
        "//app/interface/main/creative/dao/tag:go_default_library",
        "//app/interface/main/creative/dao/template:go_default_library",
        "//app/interface/main/creative/model/activity:go_default_library",
        "//app/interface/main/creative/model/appeal:go_default_library",
        "//app/interface/main/creative/model/archive:go_default_library",
        "//app/interface/main/creative/model/game:go_default_library",
        "//app/interface/main/creative/model/order:go_default_library",
        "//app/interface/main/creative/model/search:go_default_library",
        "//app/interface/main/creative/model/tag:go_default_library",
        "//app/interface/main/creative/service:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
