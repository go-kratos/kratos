package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "block.go",
        "exp.go",
        "member.go",
        "moral.go",
        "property_review.go",
        "realname.go",
    ],
    importpath = "go-common/app/service/main/member/server/grpc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/member/api:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/member/service:go_default_library",
        "//app/service/main/member/service/block:go_default_library",
        "//library/net/rpc/warden:go_default_library",
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
