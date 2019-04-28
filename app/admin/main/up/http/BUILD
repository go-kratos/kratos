package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["http_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "command.go",
        "http.go",
        "http_base.go",
        "signup.go",
        "up_data.go",
        "up_trend.go",
        "upcrm.go",
        "upload.go",
    ],
    importmap = "go-common/app/admin/main/up/http",
    importpath = "go-common/app/admin/main/up/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//app/admin/main/up/dao/upcrm:go_default_library",
        "//app/admin/main/up/model:go_default_library",
        "//app/admin/main/up/model/datamodel:go_default_library",
        "//app/admin/main/up/model/signmodel:go_default_library",
        "//app/admin/main/up/model/upcrmmodel:go_default_library",
        "//app/admin/main/up/service:go_default_library",
        "//app/admin/main/up/util:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//vendor/github.com/siddontang/go-mysql/mysql:go_default_library",
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
