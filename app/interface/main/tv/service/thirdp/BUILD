package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "dangbei_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/tv/conf:go_default_library",
        "//app/interface/main/tv/dao/thirdp:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "dangbei.go",
        "mango.go",
        "mango_recom.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/tv/service/thirdp",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/tv/conf:go_default_library",
        "//app/interface/main/tv/dao/archive:go_default_library",
        "//app/interface/main/tv/dao/cms:go_default_library",
        "//app/interface/main/tv/dao/thirdp:go_default_library",
        "//app/interface/main/tv/model:go_default_library",
        "//app/interface/main/tv/model/thirdp:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/cache:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
