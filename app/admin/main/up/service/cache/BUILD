package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cache.go",
        "tag.go",
        "video_types.go",
    ],
    importpath = "go-common/app/admin/main/up/service/cache",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//app/interface/main/creative/model/tag:go_default_library",
        "//app/service/main/videoup/model/archive:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "cache_test.go",
        "tag_test.go",
        "video_types_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//app/interface/main/creative/model/tag:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
