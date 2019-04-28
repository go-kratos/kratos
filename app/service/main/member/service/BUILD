package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "base.go",
        "exp.go",
        "member.go",
        "moral.go",
        "property_review.go",
        "realname.go",
        "service.go",
        "user_flag.go",
    ],
    importpath = "go-common/app/service/main/member/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/member/api:go_default_library",
        "//app/service/main/member/conf:go_default_library",
        "//app/service/main/member/dao:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/member/service/block:go_default_library",
        "//app/service/main/member/service/crypto:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/service/main/member/service/block:all-srcs",
        "//app/service/main/member/service/crypto:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "member_test.go",
        "realname_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/member/conf:go_default_library",
        "//app/service/main/member/dao:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
