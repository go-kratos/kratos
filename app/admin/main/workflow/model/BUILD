package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//app/admin/main/workflow/model/manager:all-srcs",
        "//app/admin/main/workflow/model/param:all-srcs",
        "//app/admin/main/workflow/model/search:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["challenge_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/smartystreets/goconvey/convey:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "activity.go",
        "block.go",
        "business.go",
        "business_attr.go",
        "callback.go",
        "challenge.go",
        "common.go",
        "event.go",
        "extra.go",
        "group.go",
        "message.go",
        "platform_meta.go",
        "reply.go",
        "state.go",
        "tag.go",
        "wlog.go",
    ],
    importpath = "go-common/app/admin/main/workflow/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/workflow/model/param:go_default_library",
        "//app/job/main/credit/model:go_default_library",
        "//library/time:go_default_library",
    ],
)
