package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "apply.go",
        "clientmoni.go",
        "cluster.go",
        "comment.go",
        "file.go",
        "grpc.go",
        "http.go",
        "job.go",
        "label.go",
        "order.go",
        "order_admin.go",
        "ptest.go",
        "rank.go",
        "report.go",
        "scene.go",
        "script.go",
        "tree.go",
        "user.go",
    ],
    importpath = "go-common/app/admin/ep/melloi/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/ep/melloi/conf:go_default_library",
        "//app/admin/ep/melloi/model:go_default_library",
        "//app/admin/ep/melloi/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
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
