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
        "//library/container/pool:all-srcs",
        "//library/container/queue/aqm:all-srcs",
    ],
    tags = ["automanaged"],
)
