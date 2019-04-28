package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "basic.go",
        "branch.go",
        "commit.go",
        "config.go",
        "http.go",
        "job.go",
        "member.go",
        "mr.go",
        "pipeline.go",
        "project.go",
        "runner.go",
        "task.go",
        "team.go",
        "user.go",
        "wechat.go",
    ],
    importpath = "go-common/app/admin/ep/saga/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/ep/saga/conf:go_default_library",
        "//app/admin/ep/saga/model:go_default_library",
        "//app/admin/ep/saga/service:go_default_library",
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
