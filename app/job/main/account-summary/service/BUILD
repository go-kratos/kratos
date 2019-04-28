package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "block.go",
        "initial.go",
        "member.go",
        "passport.go",
        "relation.go",
        "service.go",
        "summary.go",
        "sync.go",
    ],
    importpath = "go-common/app/job/main/account-summary/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/account-summary/conf:go_default_library",
        "//app/job/main/account-summary/dao:go_default_library",
        "//app/job/main/account-summary/model:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/member/model/block:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/time:go_default_library",
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
