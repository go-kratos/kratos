package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "http.go",
        "job.go",
        "local.go",
        "notice.go",
        "res.go",
    ],
    importpath = "go-common/app/interface/main/web-show/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/web-show/conf:go_default_library",
        "//app/interface/main/web-show/model/operation:go_default_library",
        "//app/interface/main/web-show/model/resource:go_default_library",
        "//app/interface/main/web-show/service/job:go_default_library",
        "//app/interface/main/web-show/service/operation:go_default_library",
        "//app/interface/main/web-show/service/resource:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/cache:go_default_library",
        "//library/net/http/blademaster/middleware/cache/store:go_default_library",
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
