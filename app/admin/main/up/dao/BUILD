package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "upload_test.go",
        "utils_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "upload.go",
        "utils.go",
    ],
    importmap = "go-common/app/admin/main/up/dao",
    importpath = "go-common/app/admin/main/up/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//library/log:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/admin/main/up/dao/data:all-srcs",
        "//app/admin/main/up/dao/email:all-srcs",
        "//app/admin/main/up/dao/global:all-srcs",
        "//app/admin/main/up/dao/manager:all-srcs",
        "//app/admin/main/up/dao/tag:all-srcs",
        "//app/admin/main/up/dao/upcrm:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
