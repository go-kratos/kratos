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
    name = "timemachine_proto",
    srcs = ["timemachine.proto"],
    tags = ["automanaged"],
    deps = ["@gogo_special_proto//github.com/gogo/protobuf/gogoproto"],
)

go_proto_library(
    name = "timemachine_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_proto"],
    importpath = "go-common/app/interface/main/activity/model/timemachine",
    proto = ":timemachine_proto",
    tags = ["automanaged"],
    deps = ["@com_github_gogo_protobuf//gogoproto:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = ["model.go"],
    embed = [":timemachine_go_proto"],
    importpath = "go-common/app/interface/main/activity/model/timemachine",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
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
