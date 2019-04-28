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
        "//app/interface/main/favorite/conf:go_default_library",
        "//app/interface/main/favorite/service:go_default_library",
        "//library/net/http/blademaster:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "http.go",
        "topic.go",
        "video.go",
    ],
    importpath = "go-common/app/interface/main/favorite/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/favorite/conf:go_default_library",
        "//app/interface/main/favorite/model:go_default_library",
        "//app/interface/main/favorite/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/anticheat:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/antispam:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/supervisor:go_default_library",
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
