package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["monkey_dao.go"],
    importpath = "go-common/app/admin/main/apm/dao/mock",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/apm/dao:go_default_library",
        "//app/admin/main/apm/model/canal:go_default_library",
        "//app/admin/main/apm/model/ecode:go_default_library",
        "//app/admin/main/apm/model/need:go_default_library",
        "//app/admin/main/apm/model/pprof:go_default_library",
        "//app/admin/main/apm/model/ut:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
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
