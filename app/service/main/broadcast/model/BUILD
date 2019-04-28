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
    name = "model_proto",
    srcs = ["model.proto"],
    tags = ["automanaged"],
)

go_proto_library(
    name = "model_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_proto"],
    importpath = "go-common/app/service/main/broadcast/model",
    proto = ":model_proto",
    tags = ["manual"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "auth.go",
        "model.go",
        "model_v1.go",
        "online.go",
        "operation.go",
        "room.go",
        "server.go",
    ],
    embed = [":model_go_proto"],
    importpath = "go-common/app/service/main/broadcast/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/broadcast/libs/bufio:go_default_library",
        "//app/service/main/broadcast/libs/bytes:go_default_library",
        "//app/service/main/broadcast/libs/encoding/binary:go_default_library",
        "//app/service/main/broadcast/libs/websocket:go_default_library",
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
