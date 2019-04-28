package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "app.go",
        "canal.go",
        "dapper.go",
        "databus.go",
        "discovery.go",
        "ecode.go",
        "http.go",
        "monitor.go",
        "need.go",
        "open.go",
        "platform.go",
        "pprof.go",
        "rank.go",
        "upload.go",
        "user.go",
        "ut.go",
        "ut_app.go",
        "ut_dash.go",
    ],
    importpath = "go-common/app/admin/main/apm/http",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/apm/conf:go_default_library",
        "//app/admin/main/apm/model/app:go_default_library",
        "//app/admin/main/apm/model/canal:go_default_library",
        "//app/admin/main/apm/model/databus:go_default_library",
        "//app/admin/main/apm/model/ecode:go_default_library",
        "//app/admin/main/apm/model/need:go_default_library",
        "//app/admin/main/apm/model/pprof:go_default_library",
        "//app/admin/main/apm/model/user:go_default_library",
        "//app/admin/main/apm/model/ut:go_default_library",
        "//app/admin/main/apm/service:go_default_library",
        "//app/tool/saga/model:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "apply_test.go",
        "canal_test.go",
        "ecode_test.go",
        "need_test.go",
        "notify_test.go",
        "process_test.go",
        "upload_test.go",
        "user_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/apm/conf:go_default_library",
        "//app/admin/main/apm/model/canal:go_default_library",
        "//app/admin/main/apm/model/need:go_default_library",
        "//app/admin/main/apm/service:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
