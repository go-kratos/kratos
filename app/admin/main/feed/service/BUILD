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
        "//app/admin/main/feed/service/bfs:all-srcs",
        "//app/admin/main/feed/service/channel:all-srcs",
        "//app/admin/main/feed/service/common:all-srcs",
        "//app/admin/main/feed/service/egg:all-srcs",
        "//app/admin/main/feed/service/pgc:all-srcs",
        "//app/admin/main/feed/service/popular:all-srcs",
        "//app/admin/main/feed/service/search:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
