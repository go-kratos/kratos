package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["recommend_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/recommend/api/grpc/v1:go_default_library",
        "//app/service/live/recommend/internal/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["recommend.go"],
    importpath = "go-common/app/service/live/recommend/internal/service/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/recommend/api/grpc/v1:go_default_library",
        "//app/service/live/recommend/internal/conf:go_default_library",
        "//app/service/live/recommend/internal/dao:go_default_library",
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
