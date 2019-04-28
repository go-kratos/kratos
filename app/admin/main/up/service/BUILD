package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "content_ext_test.go",
        "service_test.go",
        "upload_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "content_ext.go",
        "service.go",
        "upload.go",
    ],
    importmap = "go-common/app/admin/main/up/service",
    importpath = "go-common/app/admin/main/up/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//app/admin/main/up/dao:go_default_library",
        "//app/admin/main/up/dao/global:go_default_library",
        "//app/admin/main/up/dao/manager:go_default_library",
        "//app/admin/main/up/service/upcrmservice:go_default_library",
        "//app/admin/main/up/util:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/stat/prom:go_default_library",
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
        "//app/admin/main/up/service/cache:all-srcs",
        "//app/admin/main/up/service/data:all-srcs",
        "//app/admin/main/up/service/upcrmservice:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
