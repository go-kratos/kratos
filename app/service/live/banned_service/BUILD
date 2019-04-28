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
        "//app/service/live/banned_service/api/liverpc:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
