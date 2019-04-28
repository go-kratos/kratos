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
    name = "v1_proto",
    srcs = [
        "AppIndex.proto",
        "Area.proto",
        "Room.proto",
        "RoomEx.proto",
        "RoomMng.proto",
        "RoomPendant.proto",
        "RoomRecommend.proto",
        "Skin.proto",
    ],
    tags = ["automanaged"],
    deps = ["@gogo_special_proto//github.com/gogo/protobuf/gogoproto"],
)

go_proto_library(
    name = "v1_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_grpc"],
    importpath = "go-common/app/service/live/room/api/liverpc/v1",
    proto = ":v1_proto",
    tags = ["automanaged"],
    deps = ["@com_github_gogo_protobuf//gogoproto:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "AppIndex.liverpc.go",
        "Area.liverpc.go",
        "Room.liverpc.go",
        "RoomEx.liverpc.go",
        "RoomMng.liverpc.go",
        "RoomPendant.liverpc.go",
        "RoomRecommend.liverpc.go",
        "Skin.liverpc.go",
    ],
    embed = [":v1_go_proto"],
    importpath = "go-common/app/service/live/room/api/liverpc/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/net/rpc/liverpc:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
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
