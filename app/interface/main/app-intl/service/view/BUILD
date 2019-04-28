package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "access.go",
        "contain.go",
        "infoc.go",
        "service.go",
        "view.go",
    ],
    importpath = "go-common/app/interface/main/app-intl/service/view",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-intl/conf:go_default_library",
        "//app/interface/main/app-intl/dao/account:go_default_library",
        "//app/interface/main/app-intl/dao/archive:go_default_library",
        "//app/interface/main/app-intl/dao/assist:go_default_library",
        "//app/interface/main/app-intl/dao/audio:go_default_library",
        "//app/interface/main/app-intl/dao/bangumi:go_default_library",
        "//app/interface/main/app-intl/dao/coin:go_default_library",
        "//app/interface/main/app-intl/dao/dm:go_default_library",
        "//app/interface/main/app-intl/dao/favorite:go_default_library",
        "//app/interface/main/app-intl/dao/location:go_default_library",
        "//app/interface/main/app-intl/dao/manager:go_default_library",
        "//app/interface/main/app-intl/dao/region:go_default_library",
        "//app/interface/main/app-intl/dao/relation:go_default_library",
        "//app/interface/main/app-intl/dao/resource:go_default_library",
        "//app/interface/main/app-intl/dao/tag:go_default_library",
        "//app/interface/main/app-intl/dao/thumbup:go_default_library",
        "//app/interface/main/app-intl/dao/vip:go_default_library",
        "//app/interface/main/app-intl/model:go_default_library",
        "//app/interface/main/app-intl/model/bangumi:go_default_library",
        "//app/interface/main/app-intl/model/manager:go_default_library",
        "//app/interface/main/app-intl/model/region:go_default_library",
        "//app/interface/main/app-intl/model/view:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/location/model:go_default_library",
        "//app/service/main/thumbup/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/text/translate/chinese:go_default_library",
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
