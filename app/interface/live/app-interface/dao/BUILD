package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "api.go",
        "dao.go",
    ],
    importpath = "go-common/app/interface/live/app-interface/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/app-interface/conf:go_default_library",
        "//app/service/live/av/api/liverpc:go_default_library",
        "//app/service/live/fans_medal/api/liverpc:go_default_library",
        "//app/service/live/live_data/api/liverpc:go_default_library",
        "//app/service/live/live_user/api/liverpc:go_default_library",
        "//app/service/live/rankdb/api/liverpc:go_default_library",
        "//app/service/live/relation/api/liverpc:go_default_library",
        "//app/service/live/room/api/liverpc:go_default_library",
        "//app/service/live/room_ex/api/liverpc:go_default_library",
        "//app/service/live/third_api/bvc:go_default_library",
        "//app/service/live/third_api/skyhorse:go_default_library",
        "//app/service/live/userext/api/liverpc:go_default_library",
        "//app/service/live/xuser/api/grpc/v1:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
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
        "//app/interface/live/app-interface/dao/account:all-srcs",
        "//app/interface/live/app-interface/dao/av:all-srcs",
        "//app/interface/live/app-interface/dao/fans_medal:all-srcs",
        "//app/interface/live/app-interface/dao/live_data:all-srcs",
        "//app/interface/live/app-interface/dao/live_user:all-srcs",
        "//app/interface/live/app-interface/dao/rankdb:all-srcs",
        "//app/interface/live/app-interface/dao/relation:all-srcs",
        "//app/interface/live/app-interface/dao/room:all-srcs",
        "//app/interface/live/app-interface/dao/room_ex:all-srcs",
        "//app/interface/live/app-interface/dao/user_ext:all-srcs",
        "//app/interface/live/app-interface/dao/xuser:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
