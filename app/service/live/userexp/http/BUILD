load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "exp.go",
        "http.go",
    ],
    importpath = "go-common/app/service/live/userexp/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/userexp/conf:go_default_library",
        "//app/service/live/userexp/model:go_default_library",
        "//app/service/live/userexp/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/xstr:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["http_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/live/userexp/conf:go_default_library",
        "//app/service/live/userexp/model:go_default_library",
        "//app/service/live/userexp/service:go_default_library",
        "//library/net/http/blademaster:go_default_library",
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
