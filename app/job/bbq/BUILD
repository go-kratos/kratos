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
        "//app/job/bbq/cms:all-srcs",
        "//app/job/bbq/comment:all-srcs",
        "//app/job/bbq/recall:all-srcs",
        "//app/job/bbq/video:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
