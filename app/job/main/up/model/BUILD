package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "model.go",
        "reply.go",
        "request.go",
    ],
    importpath = "go-common/app/job/main/up/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
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
        "//app/job/main/up/model/archivemodel:all-srcs",
        "//app/job/main/up/model/signmodel:all-srcs",
        "//app/job/main/up/model/upcrmmodel:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
