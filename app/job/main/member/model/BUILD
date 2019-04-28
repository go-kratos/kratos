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

proto_library(
    name = "model_proto",
    srcs = ["member.proto"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["@gogo_special_proto//github.com/gogo/protobuf/gogoproto"],
)

go_proto_library(
    name = "model_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_proto"],
    importpath = "go-common/app/job/main/member/model",
    proto = ":model_proto",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["model_test.go"],
    embed = [
        ":go_default_library",
    ],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/time:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "base.go",
        "model.go",
        "moral.go",
        "notify.go",
        "official.go",
        "old_account.go",
        "realname.go",
    ],
    embed = [":model_go_proto"],
    importpath = "go-common/app/job/main/member/model",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/member/conf:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/satori/go.uuid:go_default_library",
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
        "//app/job/main/member/model/block:all-srcs",
        "//app/job/main/member/model/queue:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
