package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)
load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

proto_library(
    name = "model_proto",
    srcs = ["dm.proto"],
    tags = ["automanaged"],
    deps = ["@gogo_special_proto//github.com/gogo/protobuf/gogoproto"],
)

go_proto_library(
    name = "model_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_proto"],
    importpath = "go-common/app/job/main/dm2/model",
    proto = ":model_proto",
    tags = ["manual"],
    deps = [
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "action.go",
        "bnj.go",
        "dm.go",
        "dm_seg.go",
        "dm_special.go",
        "rank_list.go",
        "report.go",
        "subject.go",
        "subtitle.go",
        "task.go",
        "transfer.go",
        "videoup.go",
    ],
    embed = [":model_go_proto"],
    importpath = "go-common/app/job/main/dm2/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/time:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/job/main/dm2/model/oplog:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
