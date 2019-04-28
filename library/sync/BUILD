package(default_visibility = ["//visibility:public"])

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
        "//library/sync/errgroup:all-srcs",
        "//library/sync/errgroup.v2:all-srcs",
        "//library/sync/pipeline:all-srcs",
    ],
    tags = ["automanaged"],
)
