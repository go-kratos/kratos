package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "ip_old_test.go",
        "ip_test.go",
        "service_test.go",
        "zlimit_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/location/conf:go_default_library",
        "//library/cache/redis:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "anonymous.go",
        "binary.go",
        "ip.go",
        "ip_old.go",
        "service.go",
        "zlimit.go",
    ],
    importpath = "go-common/app/service/main/location/service",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/location/conf:go_default_library",
        "//app/service/main/location/dao:go_default_library",
        "//app/service/main/location/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/ip:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/ipipdotnet/ipdb-go:go_default_library",
        "//vendor/github.com/oschwald/maxminddb-golang:go_default_library",
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
