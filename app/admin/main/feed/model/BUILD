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
        "//app/admin/main/feed/model/channel:all-srcs",
        "//app/admin/main/feed/model/common:all-srcs",
        "//app/admin/main/feed/model/egg:all-srcs",
        "//app/admin/main/feed/model/search:all-srcs",
        "//app/admin/main/feed/model/show:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
