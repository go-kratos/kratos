load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)
load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cloud.go",
        "rpc.go",
    ],
    embed = [":model_go_proto"],
    importpath = "go-common/app/service/main/identify-game/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
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

proto_library(
    name = "model_proto",
    srcs = ["token.proto"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["@gogo_special_proto//github.com/gogo/protobuf/gogoproto"],
)

go_proto_library(
    name = "model_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_proto"],
    importpath = "go-common/app/service/main/identify-game/model",
    proto = ":model_proto",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["@com_github_gogo_protobuf//gogoproto:go_default_library"],
)
