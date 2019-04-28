package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "contribute_test.go",
        "service_test.go",
        "space_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-interface/conf:go_default_library",
        "//app/interface/main/app-interface/model:go_default_library",
        "//app/interface/main/app-interface/model/space:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "contribute.go",
        "service.go",
        "space.go",
    ],
    importpath = "go-common/app/interface/main/app-interface/service/space",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-interface/conf:go_default_library",
        "//app/interface/main/app-interface/dao/account:go_default_library",
        "//app/interface/main/app-interface/dao/archive:go_default_library",
        "//app/interface/main/app-interface/dao/article:go_default_library",
        "//app/interface/main/app-interface/dao/audio:go_default_library",
        "//app/interface/main/app-interface/dao/bangumi:go_default_library",
        "//app/interface/main/app-interface/dao/bplus:go_default_library",
        "//app/interface/main/app-interface/dao/coin:go_default_library",
        "//app/interface/main/app-interface/dao/community:go_default_library",
        "//app/interface/main/app-interface/dao/elec:go_default_library",
        "//app/interface/main/app-interface/dao/favorite:go_default_library",
        "//app/interface/main/app-interface/dao/live:go_default_library",
        "//app/interface/main/app-interface/dao/member:go_default_library",
        "//app/interface/main/app-interface/dao/pay:go_default_library",
        "//app/interface/main/app-interface/dao/relation:go_default_library",
        "//app/interface/main/app-interface/dao/search:go_default_library",
        "//app/interface/main/app-interface/dao/shop:go_default_library",
        "//app/interface/main/app-interface/dao/space:go_default_library",
        "//app/interface/main/app-interface/dao/tag:go_default_library",
        "//app/interface/main/app-interface/dao/thumbup:go_default_library",
        "//app/interface/main/app-interface/model:go_default_library",
        "//app/interface/main/app-interface/model/audio:go_default_library",
        "//app/interface/main/app-interface/model/bangumi:go_default_library",
        "//app/interface/main/app-interface/model/bplus:go_default_library",
        "//app/interface/main/app-interface/model/community:go_default_library",
        "//app/interface/main/app-interface/model/favorite:go_default_library",
        "//app/interface/main/app-interface/model/shop:go_default_library",
        "//app/interface/main/app-interface/model/space:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/thumbup/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
