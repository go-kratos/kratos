package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "channel.go",
        "full.go",
        "http.go",
        "pgc.go",
        "recruit.go",
        "share.go",
        "ugc.go",
        "wechat.go",
    ],
    importpath = "go-common/app/interface/main/web-goblin/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/web-goblin/conf:go_default_library",
        "//app/interface/main/web-goblin/model/web:go_default_library",
        "//app/interface/main/web-goblin/service/share:go_default_library",
        "//app/interface/main/web-goblin/service/web:go_default_library",
        "//app/interface/main/web-goblin/service/wechat:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
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
