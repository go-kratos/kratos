package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "full_test.go",
        "pgc_test.go",
        "recruit_test.go",
        "service_test.go",
        "ugc_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/web-goblin/conf:go_default_library",
        "//app/interface/main/web-goblin/model/web:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "channel.go",
        "full.go",
        "pgc.go",
        "recruit.go",
        "service.go",
        "ugc.go",
    ],
    importpath = "go-common/app/interface/main/web-goblin/service/web",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/tag/model:go_default_library",
        "//app/interface/main/tag/rpc/client:go_default_library",
        "//app/interface/main/web-goblin/conf:go_default_library",
        "//app/interface/main/web-goblin/dao/web:go_default_library",
        "//app/interface/main/web-goblin/model/web:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/xstr:go_default_library",
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
