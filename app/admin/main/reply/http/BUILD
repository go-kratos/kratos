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
        "//app/admin/main/reply/conf:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "action.go",
        "business.go",
        "emoji.go",
        "emoji_package.go",
        "http.go",
        "log.go",
        "monitor.go",
        "notice.go",
        "reply.go",
        "report.go",
        "subject.go",
    ],
    importpath = "go-common/app/admin/main/reply/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/reply/conf:go_default_library",
        "//app/admin/main/reply/model:go_default_library",
        "//app/admin/main/reply/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/ip:go_default_library",
        "//library/time:go_default_library",
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
