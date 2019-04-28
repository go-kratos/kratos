package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["service_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/report-click/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "err_report.go",
        "ip.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/report-click/service",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/history/model:go_default_library",
        "//app/interface/main/history/rpc/client:go_default_library",
        "//app/interface/main/report-click/conf:go_default_library",
        "//app/interface/main/report-click/dao:go_default_library",
        "//app/interface/main/report-click/model:go_default_library",
        "//app/interface/main/report-click/service/crypto/aes:go_default_library",
        "//app/interface/main/report-click/service/crypto/padding:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/account/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/interface/main/report-click/service/crypto/aes:all-srcs",
        "//app/interface/main/report-click/service/crypto/cipher:all-srcs",
        "//app/interface/main/report-click/service/crypto/padding:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
