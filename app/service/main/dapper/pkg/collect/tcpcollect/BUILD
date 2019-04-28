package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["server_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/dapper/conf:go_default_library",
        "//app/service/main/dapper/model:go_default_library",
        "//app/service/main/dapper/pkg/process:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["server.go"],
    importpath = "go-common/app/service/main/dapper/pkg/collect/tcpcollect",
    tags = ["manual"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/dapper/conf:go_default_library",
        "//app/service/main/dapper/model:go_default_library",
        "//app/service/main/dapper/pkg/process:go_default_library",
        "//library/log:go_default_library",
        "//library/net/trace/proto:go_default_library",
        "//library/stat/counter:go_default_library",
        "//library/stat/prom:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@io_bazel_rules_go//proto/wkt:duration_go_proto",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
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
