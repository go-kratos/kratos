package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
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
        "//app/admin/main/apm/model/app:all-srcs",
        "//app/admin/main/apm/model/canal:all-srcs",
        "//app/admin/main/apm/model/databus:all-srcs",
        "//app/admin/main/apm/model/discovery:all-srcs",
        "//app/admin/main/apm/model/ecode:all-srcs",
        "//app/admin/main/apm/model/log:all-srcs",
        "//app/admin/main/apm/model/monitor:all-srcs",
        "//app/admin/main/apm/model/need:all-srcs",
        "//app/admin/main/apm/model/pprof:all-srcs",
        "//app/admin/main/apm/model/tree:all-srcs",
        "//app/admin/main/apm/model/user:all-srcs",
        "//app/admin/main/apm/model/ut:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
