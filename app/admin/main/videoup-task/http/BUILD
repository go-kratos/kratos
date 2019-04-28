package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "consumer.go",
        "http.go",
        "prehandler.go",
        "qa_video.go",
        "report.go",
        "review.go",
        "task.go",
        "weight.go",
    ],
    importpath = "go-common/app/admin/main/videoup-task/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/videoup-task/conf:go_default_library",
        "//app/admin/main/videoup-task/model:go_default_library",
        "//app/admin/main/videoup-task/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
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
