package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "app_info_test.go",
        "notice_test.go",
        "remotelogin_test.go",
        "reply_test.go",
        "service_test.go",
        "web_info_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/account/conf:go_default_library",
        "//app/interface/main/account/model:go_default_library",
        "//app/service/main/secure/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "app_info.go",
        "medal.go",
        "member.go",
        "notice.go",
        "remotelogin.go",
        "reply.go",
        "service.go",
        "sudo.go",
        "web_info.go",
    ],
    importpath = "go-common/app/interface/main/account/service/member",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/account/conf:go_default_library",
        "//app/interface/main/account/dao/account:go_default_library",
        "//app/interface/main/account/dao/passport:go_default_library",
        "//app/interface/main/account/dao/reply:go_default_library",
        "//app/interface/main/account/model:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/interface/openplatform/article/rpc/client:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/account/rpc/client:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/coin/api/gorpc:go_default_library",
        "//app/service/main/coin/model:go_default_library",
        "//app/service/main/filter/model/rpc:go_default_library",
        "//app/service/main/filter/rpc/client:go_default_library",
        "//app/service/main/location/model:go_default_library",
        "//app/service/main/location/rpc/client:go_default_library",
        "//app/service/main/member/api/gorpc:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/passport/model:go_default_library",
        "//app/service/main/passport/rpc/client:go_default_library",
        "//app/service/main/secure/model:go_default_library",
        "//app/service/main/secure/rpc/client:go_default_library",
        "//app/service/main/up/api/gorpc:go_default_library",
        "//app/service/main/usersuit/model:go_default_library",
        "//app/service/main/usersuit/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/queue/databus/report:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
