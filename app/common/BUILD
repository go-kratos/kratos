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
        "//app/common/live:all-srcs",
        "//app/common/openplatform:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
