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
    srcs = ["message.proto"],
    tags = ["automanaged"],
    deps = ["@gogo_special_proto//github.com/gogo/protobuf/gogoproto"],
)

go_proto_library(
    name = "model_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_proto"],
    importpath = "go-common/app/service/main/up/model",
    proto = ":model_proto",
    tags = ["automanaged"],
    deps = [
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "monitor.go",
        "msg.go",
        "rpc.go",
        "sign_up.go",
        "up.go",
        "up_card.go",
        "up_group.go",
        "up_info.go",
        "up_special.go",
        "up_switch.go",
    ],
    embed = [":model_go_proto"],
    importpath = "go-common/app/service/main/up/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/time:go_default_library",
        "//vendor/github.com/siddontang/go-mysql/mysql:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/service/main/up/model/data:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
