package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "fansMedal.go",
        "mainAccount.go",
        "mainFilter.go",
        "relation.go",
        "resource.go",
        "room.go",
        "user.go",
    ],
    importpath = "go-common/app/interface/live/app-ucenter/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/app-ucenter/conf:go_default_library",
        "//app/service/live/fans_medal/api/liverpc:go_default_library",
        "//app/service/live/fans_medal/api/liverpc/v1:go_default_library",
        "//app/service/live/relation/api/liverpc:go_default_library",
        "//app/service/live/relation/api/liverpc/v1:go_default_library",
        "//app/service/live/resource/api/grpc/v1:go_default_library",
        "//app/service/live/room/api/liverpc:go_default_library",
        "//app/service/live/room/api/liverpc/v1:go_default_library",
        "//app/service/live/user/api/liverpc:go_default_library",
        "//app/service/live/user/api/liverpc/v3:go_default_library",
        "//app/service/main/member/api:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
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
