load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "add_ver_test.go",
        "publish_test.go",
        "service_test.go",
        "upbfs_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/appstatic/conf:go_default_library",
        "//app/admin/main/appstatic/model:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "add_ver.go",
        "publish.go",
        "push.go",
        "service.go",
        "upbfs.go",
    ],
    importpath = "go-common/app/admin/main/appstatic/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/appstatic/conf:go_default_library",
        "//app/admin/main/appstatic/dao:go_default_library",
        "//app/admin/main/appstatic/model:go_default_library",
        "//library/log:go_default_library",
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
