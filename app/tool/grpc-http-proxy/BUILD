package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
)

go_library(
    name = "go_default_library",
    srcs = [
        "grpc.go",
        "http.go",
        "main.go",
    ],
    importpath = "go-common/app/tool/grpc-http-proxy",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/naming/discovery:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//jsonpb:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//encoding:go_default_library",
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

go_binary(
    name = "grpc-http-proxy",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)
