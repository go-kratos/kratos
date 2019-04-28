package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "card.go",
        "equip.go",
        "group.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/card/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/card/conf:go_default_library",
        "//app/service/main/card/dao:go_default_library",
        "//app/service/main/card/model:go_default_library",
        "//app/service/main/vip/model:go_default_library",
        "//app/service/main/vip/rpc/client:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "card_test.go",
        "equip_test.go",
        "group_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/card/conf:go_default_library",
        "//app/service/main/card/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
