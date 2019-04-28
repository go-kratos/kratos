package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cursor.go",
        "http.go",
        "local.go",
        "reply.go",
        "reply_admin.go",
        "reply_record.go",
        "reply_report.go",
        "xreply.go",
    ],
    importpath = "go-common/app/interface/main/reply/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/reply/conf:go_default_library",
        "//app/interface/main/reply/model/adminlog:go_default_library",
        "//app/interface/main/reply/model/reply:go_default_library",
        "//app/interface/main/reply/model/xreply:go_default_library",
        "//app/interface/main/reply/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/metadata:go_default_library",
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
