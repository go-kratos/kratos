package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)
load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "addit.go",
        "archive.go",
        "pgc.go",
        "report_result.go",
        "rpc.go",
        "rpc2.go",
    ],
    embed = [":archive_go_proto"],
    importpath = "go-common/app/service/main/archive/model/archive",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
)

proto_library(
    name = "archive_proto",
    srcs = [
        "archive.proto",
        "bvc.proto",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["@gogo_special_proto//github.com/gogo/protobuf/gogoproto"],
)

go_proto_library(
    name = "archive_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_proto"],
    importpath = "go-common/app/service/main/archive/model/archive",
    proto = ":archive_proto",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
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
