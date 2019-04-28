package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "act.go",
        "bnj.go",
        "bws.go",
        "http.go",
        "kfc.go",
        "like.go",
        "match.go",
        "sports.go",
        "timemachine.go",
    ],
    importpath = "go-common/app/interface/main/activity/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/activity/conf:go_default_library",
        "//app/interface/main/activity/model/bws:go_default_library",
        "//app/interface/main/activity/model/like:go_default_library",
        "//app/interface/main/activity/model/sports:go_default_library",
        "//app/interface/main/activity/service/bnj:go_default_library",
        "//app/interface/main/activity/service/bws:go_default_library",
        "//app/interface/main/activity/service/kfc:go_default_library",
        "//app/interface/main/activity/service/like:go_default_library",
        "//app/interface/main/activity/service/sports:go_default_library",
        "//app/interface/main/activity/service/timemachine:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/metadata:go_default_library",
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
