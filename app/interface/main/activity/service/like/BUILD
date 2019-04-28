package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "act_test.go",
        "like_test.go",
        "match_test.go",
        "missiongroup_test.go",
        "service_test.go",
        "subject_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/activity/conf:go_default_library",
        "//app/interface/main/activity/model/like:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "act.go",
        "bnj.go",
        "like.go",
        "match.go",
        "missiongroup.go",
        "service.go",
        "subject.go",
    ],
    importpath = "go-common/app/interface/main/activity/service/like",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/activity/conf:go_default_library",
        "//app/interface/main/activity/dao/bnj:go_default_library",
        "//app/interface/main/activity/dao/like:go_default_library",
        "//app/interface/main/activity/model/bnj:go_default_library",
        "//app/interface/main/activity/model/like:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/interface/main/tag/rpc/client:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/coin/api/gorpc:go_default_library",
        "//app/service/main/coin/model:go_default_library",
        "//app/service/main/spy/model:go_default_library",
        "//app/service/main/spy/rpc/client:go_default_library",
        "//app/service/main/thumbup/model:go_default_library",
        "//app/service/main/thumbup/rpc/client:go_default_library",
        "//app/service/main/usersuit/model:go_default_library",
        "//app/service/main/usersuit/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
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
