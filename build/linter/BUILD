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
        "//build/linter/internal:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
