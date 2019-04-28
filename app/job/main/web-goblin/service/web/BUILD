package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "service.go",
        "ugc.go",
    ],
    importpath = "go-common/app/job/main/web-goblin/service/web",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/web-goblin/conf:go_default_library",
        "//app/job/main/web-goblin/dao/web:go_default_library",
        "//app/job/main/web-goblin/model/web:go_default_library",
        "//app/service/main/dynamic/model:go_default_library",
        "//app/service/main/dynamic/rpc/client:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
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
        "service_test.go",
        "ugc_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/job/main/web-goblin/conf:go_default_library",
        "//app/job/main/web-goblin/model/web:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
