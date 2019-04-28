package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["document_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/kvo/conf:go_default_library",
        "//app/interface/main/kvo/model/module:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/container/pool:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/net/netutil/breaker:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "document.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/kvo/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/kvo/conf:go_default_library",
        "//app/interface/main/kvo/dao:go_default_library",
        "//app/interface/main/kvo/model:go_default_library",
        "//app/interface/main/kvo/model/module:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
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
