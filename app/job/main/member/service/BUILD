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
        "//app/job/main/member/conf:go_default_library",
        "//app/job/main/member/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "base.go",
        "cache_delay.go",
        "check.go",
        "exp.go",
        "fixer.go",
        "member.go",
        "moral.go",
        "realname.go",
        "service.go",
        "subproc.go",
    ],
    importpath = "go-common/app/job/main/member/service",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/account/service/realname/crypto:go_default_library",
        "//app/job/main/member/conf:go_default_library",
        "//app/job/main/member/dao:go_default_library",
        "//app/job/main/member/model:go_default_library",
        "//app/job/main/member/model/queue:go_default_library",
        "//app/job/main/member/service/block:go_default_library",
        "//app/service/main/member/api/gorpc:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/share/model:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/ip:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/queue/databus/databusutil:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/golang.org/x/time/rate:go_default_library",
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
        "//app/job/main/member/service/block:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
