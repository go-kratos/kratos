package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "favorite_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-interface/conf:go_default_library",
        "//app/interface/main/app-interface/model/favorite:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "favorite.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/app-interface/service/favorite",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-interface/conf:go_default_library",
        "//app/interface/main/app-interface/dao/article:go_default_library",
        "//app/interface/main/app-interface/dao/audio:go_default_library",
        "//app/interface/main/app-interface/dao/bangumi:go_default_library",
        "//app/interface/main/app-interface/dao/bplus:go_default_library",
        "//app/interface/main/app-interface/dao/favorite:go_default_library",
        "//app/interface/main/app-interface/dao/mall:go_default_library",
        "//app/interface/main/app-interface/dao/sp:go_default_library",
        "//app/interface/main/app-interface/dao/ticket:go_default_library",
        "//app/interface/main/app-interface/dao/topic:go_default_library",
        "//app/interface/main/app-interface/model:go_default_library",
        "//app/interface/main/app-interface/model/favorite:go_default_library",
        "//app/service/main/favorite/api:go_default_library",
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
