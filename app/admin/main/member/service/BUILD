load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)
load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "member_test.go",
        "official_test.go",
        "review_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/member/conf:go_default_library",
        "//app/admin/main/member/model:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "member.go",
        "monitor.go",
        "official.go",
        "realname.go",
        "review.go",
        "review_audit.go",
        "service.go",
    ],
    importpath = "go-common/app/admin/main/member/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/member/conf:go_default_library",
        "//app/admin/main/member/dao:go_default_library",
        "//app/admin/main/member/model:go_default_library",
        "//app/admin/main/member/model/block:go_default_library",
        "//app/admin/main/member/model/bom:go_default_library",
        "//app/admin/main/member/service/block:go_default_library",
        "//app/interface/main/account/service/realname/crypto:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/coin/api/gorpc:go_default_library",
        "//app/service/main/coin/model:go_default_library",
        "//app/service/main/figure/rpc/client:go_default_library",
        "//app/service/main/member/api/gorpc:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/member/model/block:go_default_library",
        "//app/service/main/member/service/crypto:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//app/service/main/relation/rpc/client:go_default_library",
        "//app/service/main/spy/model:go_default_library",
        "//app/service/main/spy/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/nfnt/resize:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
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
        "//app/admin/main/member/service/block:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
