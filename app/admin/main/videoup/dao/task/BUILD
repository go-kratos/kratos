package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "weight.go",
    ],
    importpath = "go-common/app/admin/main/videoup/dao/task",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/videoup/conf:go_default_library",
        "//app/admin/main/videoup/model/archive:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/log:go_default_library",
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
        "dao_test.go",
        "weight_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/videoup/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
