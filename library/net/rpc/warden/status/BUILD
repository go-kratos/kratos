package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["status.go"],
    importpath = "go-common/library/net/rpc/warden/status",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/ecode:go_default_library",
        "//library/ecode/pb:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = ["status_test.go"],
    embed = [":go_default_library"],
    tags = ["manual"],
    deps = [
        "//library/ecode:go_default_library",
        "//library/ecode/pb:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
