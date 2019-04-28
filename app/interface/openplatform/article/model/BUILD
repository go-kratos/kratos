package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)
load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

go_test(
    name = "go_default_test",
    srcs = ["article_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/smartystreets/goconvey/convey:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "activity.go",
        "apply.go",
        "article.go",
        "author.go",
        "cards.go",
        "creation.go",
        "hotspots.go",
        "infoc.go",
        "list.go",
        "rank.go",
        "rpc.go",
        "setting.go",
        "sort.go",
    ],
    embed = [":model_go_proto"],
    importpath = "go-common/app/interface/openplatform/article/model",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/creative/model/data:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
    ],
)

proto_library(
    name = "model_proto",
    srcs = ["model.proto"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["@gogo_special_proto//github.com/gogo/protobuf/gogoproto"],
)

go_proto_library(
    name = "model_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_proto"],
    importpath = "go-common/app/interface/openplatform/article/model",
    proto = ":model_proto",
    tags = ["manual"],
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
    srcs = [
        ":package-srcs",
        "//app/interface/openplatform/article/model/search:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
