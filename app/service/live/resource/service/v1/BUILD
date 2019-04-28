package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "banner_test.go",
        "liveCheck_test.go",
        "resource_test.go",
        "splash_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/resource/api/grpc/v1:go_default_library",
        "//app/service/live/resource/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "banner.go",
        "liveCheck.go",
        "resource.go",
        "splash.go",
        "titans.go",
        "tree.go",
    ],
    importpath = "go-common/app/service/live/resource/service/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/resource/api/grpc/v1:go_default_library",
        "//app/service/live/resource/api/http/v1:go_default_library",
        "//app/service/live/resource/conf:go_default_library",
        "//app/service/live/resource/dao:go_default_library",
        "//app/service/live/resource/lrucache:go_default_library",
        "//app/service/live/resource/model:go_default_library",
        "//library/ecode:go_default_library",
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
