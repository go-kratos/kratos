load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "achieve.go",
        "addit.go",
        "attr.go",
        "audit.go",
        "following.go",
        "log.go",
        "rpc.go",
        "stat.go",
    ],
    embed = [":model_go_proto"],
    importpath = "go-common/app/service/main/relation/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
    ],
)

proto_library(
    name = "model_proto",
    srcs = [
        "model.proto",
    ],
    tags = ["manual"],
    visibility = ["//visibility:public"],
    deps = [
        "@gogo_special_proto//github.com/gogo/protobuf/gogoproto",
    ],
)

go_proto_library(
    name = "model_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogo_proto"],
    importpath = "go-common/app/service/main/relation/model",
    proto = ":model_proto",
    tags = ["manual"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_gogo_protobuf//sortkeys:go_default_library",
        "@com_github_gogo_protobuf//types:go_default_library",
        "@io_bazel_rules_go//proto/wkt:any_go_proto",
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
        "//app/service/main/relation/model/i64b:all-srcs",
        "//app/service/main/relation/model/sets:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
