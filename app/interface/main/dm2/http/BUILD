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
    deps = ["//app/interface/main/dm2/model:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "dm.go",
        "filter.go",
        "history.go",
        "http.go",
        "manage.go",
        "mask.go",
        "member.go",
        "subtitle.go",
        "thumbup.go",
        "view.go",
    ],
    importpath = "go-common/app/interface/main/dm2/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/dm2/conf:go_default_library",
        "//app/interface/main/dm2/model:go_default_library",
        "//app/interface/main/dm2/model/oplog:go_default_library",
        "//app/interface/main/dm2/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/antispam:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/ip:go_default_library",
        "//library/net/metadata:go_default_library",
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
