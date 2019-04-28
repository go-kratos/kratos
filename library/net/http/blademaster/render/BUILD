load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)
load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

proto_library(
    name = "render_proto",
    srcs = ["pb.proto"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_google_protobuf//:any_proto",
        "@gogo_special_proto//github.com/gogo/protobuf/gogoproto",
    ],
)

go_proto_library(
    name = "render_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogo_proto"],
    importpath = "go-common/library/net/http/blademaster/render",
    proto = ":render_proto",
    tags = ["manual"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "data.go",
        "json.go",
        "protobuf.go",
        "redirect.go",
        "render.go",
        "string.go",
        "xml.go",
    ],
    embed = [":render_go_proto"],
    importpath = "go-common/library/net/http/blademaster/render",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/pkg/errors:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_gogo_protobuf//types:go_default_library",
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
