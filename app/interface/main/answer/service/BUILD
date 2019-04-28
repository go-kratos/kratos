package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "answer_test.go",
        "answercheck_test.go",
        "common_test.go",
        "geetest_test.go",
        "service_test.go",
        "user_info_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/answer/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "answer.go",
        "answercheck.go",
        "common.go",
        "geetest.go",
        "rank_share.go",
        "service.go",
        "user_info.go",
    ],
    importpath = "go-common/app/interface/main/answer/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/answer/conf:go_default_library",
        "//app/interface/main/answer/dao:go_default_library",
        "//app/interface/main/answer/dao/account:go_default_library",
        "//app/interface/main/answer/dao/geetest:go_default_library",
        "//app/interface/main/answer/model:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/member/api/gorpc:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/anticheat:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/text/translate/chinese:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "@org_golang_x_net//context:go_default_library",
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
