package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "command.go",
        "contributors.go",
        "merge.go",
        "review.go",
    ],
    importpath = "go-common/app/tool/saga/service/command",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/tool/saga/conf:go_default_library",
        "//app/tool/saga/dao:go_default_library",
        "//app/tool/saga/model:go_default_library",
        "//app/tool/saga/service/gitlab:go_default_library",
        "//app/tool/saga/service/notification:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/xanzy/go-gitlab:go_default_library",
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
    srcs = [
        "command_test.go",
        "contributors_test.go",
        "review_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/tool/saga/conf:go_default_library",
        "//app/tool/saga/dao:go_default_library",
        "//app/tool/saga/model:go_default_library",
        "//app/tool/saga/service/gitlab:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
