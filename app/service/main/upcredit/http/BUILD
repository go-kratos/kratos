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
    deps = ["//app/service/main/up/conf:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "http.go",
        "logcredit.go",
    ],
    importmap = "go-common/app/service/main/upcredit/http",
    importpath = "go-common/app/service/main/upcredit/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/upcredit/conf:go_default_library",
        "//app/service/main/upcredit/model/calculator:go_default_library",
        "//app/service/main/upcredit/model/upcrmmodel:go_default_library",
        "//app/service/main/upcredit/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
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
