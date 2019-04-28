load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
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
        "//app/service/main/upcredit/model/calculator:all-srcs",
        "//app/service/main/upcredit/model/canal:all-srcs",
        "//app/service/main/upcredit/model/databus:all-srcs",
        "//app/service/main/upcredit/model/upcrmmodel:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
