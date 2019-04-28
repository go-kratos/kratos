package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "business.go",
        "challenge.go",
        "group.go",
        "http.go",
        "sign.go",
        "sobot.go",
        "tag.go",
    ],
    importpath = "go-common/app/service/main/workflow/http",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/workflow/conf:go_default_library",
        "//app/service/main/workflow/model:go_default_library",
        "//app/service/main/workflow/model/account:go_default_library",
        "//app/service/main/workflow/model/sobot:go_default_library",
        "//app/service/main/workflow/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
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
