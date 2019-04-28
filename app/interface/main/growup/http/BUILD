package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "activity.go",
        "column.go",
        "exchange.go",
        "http.go",
        "income.go",
        "join.go",
        "local.go",
        "newbie.go",
        "notice.go",
        "special_award.go",
        "web.go",
        "withdraw.go",
    ],
    importpath = "go-common/app/interface/main/growup/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/growup/conf:go_default_library",
        "//app/interface/main/growup/model:go_default_library",
        "//app/interface/main/growup/service:go_default_library",
        "//app/interface/main/growup/service/newbie:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/net/metadata:go_default_library",
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
