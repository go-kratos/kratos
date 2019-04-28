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
        "//build/linter:all-srcs",
        "//build/visible_to:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
