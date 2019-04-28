package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "audit_log.go",
        "base_model.go",
        "formatter.go",
        "mcn.go",
        "mcn_pay.go",
        "msg.go",
        "reply.go",
        "req.go",
        "request_base.go",
        "request_example.go",
        "statistics.go",
    ],
    importpath = "go-common/app/admin/main/mcn/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/mcn/model/datamodel:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = ["reply_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
)
