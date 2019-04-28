package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "achieve.go",
        "aes.go",
        "audit.go",
        "black.go",
        "follower.go",
        "following.go",
        "infoc.go",
        "monitor.go",
        "prompt.go",
        "relation.go",
        "service.go",
        "stat.go",
        "tag.go",
        "whisper.go",
    ],
    importpath = "go-common/app/service/main/relation/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/member/api/gorpc:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/member/model/block:go_default_library",
        "//app/service/main/relation/conf:go_default_library",
        "//app/service/main/relation/dao:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//app/service/main/relation/model/sets:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/stat/prom:go_default_library",
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
