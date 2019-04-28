package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "call.go",
        "dao.go",
        "lic.go",
        "xml.go",
    ],
    importpath = "go-common/app/job/main/tv/dao/lic",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/tv/conf:go_default_library",
        "//app/job/main/tv/model/pgc:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "lic_test.go",
        "xml_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/tv/conf:go_default_library",
        "//app/job/main/tv/model/pgc:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
