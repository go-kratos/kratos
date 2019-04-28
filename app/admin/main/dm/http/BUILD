package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "adv_test.go",
        "filter_test.go",
        "http_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//app/admin/main/dm/conf:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "adv.go",
        "dm.go",
        "filter.go",
        "http.go",
        "mask.go",
        "monitor.go",
        "report.go",
        "shield.go",
        "subject.go",
        "subtitle.go",
        "task.go",
        "transfer.go",
    ],
    importpath = "go-common/app/admin/main/dm/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/dm/conf:go_default_library",
        "//app/admin/main/dm/model:go_default_library",
        "//app/admin/main/dm/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/xstr:go_default_library",
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
