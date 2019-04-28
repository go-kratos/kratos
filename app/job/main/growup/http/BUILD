package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "blacklist.go",
        "bubble.go",
        "cheat.go",
        "combinemails.go",
        "http.go",
        "tag.go",
        "tagincome.go",
        "task_status.go",
        "web.go",
    ],
    importpath = "go-common/app/job/main/growup/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/service:go_default_library",
        "//app/job/main/growup/service/charge:go_default_library",
        "//app/job/main/growup/service/ctrl:go_default_library",
        "//app/job/main/growup/service/income:go_default_library",
        "//app/job/main/growup/service/tag:go_default_library",
        "//library/log:go_default_library",
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
