package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "relation_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-interface/conf:go_default_library",
        "//app/interface/main/app-interface/model/relation:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "relation.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/app-interface/service/relation",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-interface/conf:go_default_library",
        "//app/interface/main/app-interface/dao/account:go_default_library",
        "//app/interface/main/app-interface/dao/live:go_default_library",
        "//app/interface/main/app-interface/dao/relation:go_default_library",
        "//app/interface/main/app-interface/model/live:go_default_library",
        "//app/interface/main/app-interface/model/relation:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//library/ecode:go_default_library",
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
