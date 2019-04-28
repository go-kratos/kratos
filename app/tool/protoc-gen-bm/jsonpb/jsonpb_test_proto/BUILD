load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

proto_library(
    name = "jsonpb_proto",
    srcs = [
        "more_test_objects.proto",
        "test_objects.proto",
    ],
    tags = ["manual"],
    deps = [
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:duration_proto",
        "@com_google_protobuf//:struct_proto",
        "@com_google_protobuf//:timestamp_proto",
        "@com_google_protobuf//:wrappers_proto",
    ],
)

go_proto_library(
    name = "jsonpb_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "go-common/app/tool/protoc-gen-bm/jsonpb/jsonpb_test_proto",
    proto = ":jsonpb_proto",
    tags = ["manual"],
    deps = [
        "@com_github_golang_protobuf//ptypes/struct:go_default_library",
        "@com_github_golang_protobuf//ptypes/wrappers:go_default_library",
        "@io_bazel_rules_go//proto/wkt:any_go_proto",
        "@io_bazel_rules_go//proto/wkt:duration_go_proto",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [],
    embed = [":jsonpb_go_proto"],
    importpath = "go-common/app/tool/protoc-gen-bm/jsonpb/jsonpb_test_proto",
    tags = ["manual"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//ptypes/struct:go_default_library",
        "@com_github_golang_protobuf//ptypes/wrappers:go_default_library",
        "@io_bazel_rules_go//proto/wkt:any_go_proto",
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
