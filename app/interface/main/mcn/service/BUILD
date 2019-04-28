package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "account.go",
        "api.go",
        "cmd.go",
        "creative.go",
        "data.go",
        "mcn.go",
        "msg.go",
        "rank.go",
        "service.go",
        "uniquecheck.go",
        "upload.go",
    ],
    importpath = "go-common/app/interface/main/mcn/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/mcn/model:go_default_library",
        "//app/interface/main/mcn/conf:go_default_library",
        "//app/interface/main/mcn/dao/bfs:go_default_library",
        "//app/interface/main/mcn/dao/cache:go_default_library",
        "//app/interface/main/mcn/dao/datadao:go_default_library",
        "//app/interface/main/mcn/dao/global:go_default_library",
        "//app/interface/main/mcn/dao/mcndao:go_default_library",
        "//app/interface/main/mcn/dao/msg:go_default_library",
        "//app/interface/main/mcn/model:go_default_library",
        "//app/interface/main/mcn/model/datamodel:go_default_library",
        "//app/interface/main/mcn/model/mcnmodel:go_default_library",
        "//app/interface/main/mcn/tool/worker:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/member/api:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/member/model/block:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/bluele/gcache:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
        "account_test.go",
        "cmd_test.go",
        "mcn_test.go",
        "msg_test.go",
        "rank_test.go",
        "service_test.go",
        "uniquecheck_test.go",
        "upload_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/mcn/model:go_default_library",
        "//app/interface/main/mcn/conf:go_default_library",
        "//app/interface/main/mcn/model:go_default_library",
        "//app/interface/main/mcn/model/mcnmodel:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
