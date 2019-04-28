package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "arctag_action_test.go",
        "arctag_changeapi_test.go",
        "arctag_user_test.go",
        "channel_test.go",
        "restag_test.go",
        "service_test.go",
        "sub_sort_test.go",
        "sub_test.go",
        "tag_test.go",
        "tagarc_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/tag/conf:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/cache/redis:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "arctag_action.go",
        "arctag_changeapi.go",
        "arctag_check.go",
        "arctag_user.go",
        "channel.go",
        "http.go",
        "new_arc.go",
        "restag.go",
        "rpc.go",
        "rpc_channel.go",
        "rpc_tag.go",
        "service.go",
        "sub.go",
        "sub_sort.go",
        "tag.go",
        "tagarc.go",
    ],
    importpath = "go-common/app/interface/main/tag/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/tag/conf:go_default_library",
        "//app/interface/main/tag/dao:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/figure/model:go_default_library",
        "//app/service/main/figure/rpc/client:go_default_library",
        "//app/service/main/tag/api:go_default_library",
        "//app/service/main/tag/model:go_default_library",
        "//app/service/main/tag/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/golang.org/x/sync/singleflight:go_default_library",
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
