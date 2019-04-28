load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cloud.go",
        "identify.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/identify-game/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/identify-game/api/grpc/v1:go_default_library",
        "//app/service/main/identify-game/conf:go_default_library",
        "//app/service/main/identify-game/dao:go_default_library",
        "//app/service/main/identify-game/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/stat:go_default_library",
        "//library/stat/prom:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "cloud_test.go",
        "identify_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/identify-game/conf:go_default_library",
        "//app/service/main/identify-game/model:go_default_library",
        "//library/ecode:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
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
