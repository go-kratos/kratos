package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
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
        "archive.go",
        "attachment.go",
        "business.go",
        "callback.go",
        "challenge.go",
        "control.go",
        "event.go",
        "group.go",
        "log.go",
        "tag.go",
    ],
    importpath = "go-common/app/service/main/workflow/model",
    tags = ["automanaged"],
    deps = ["//library/time:go_default_library"],
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
        "//app/service/main/workflow/model/account:all-srcs",
        "//app/service/main/workflow/model/sobot:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
