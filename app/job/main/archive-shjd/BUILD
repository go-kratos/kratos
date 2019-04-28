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
        "//app/job/main/archive-shjd/cmd:all-srcs",
        "//app/job/main/archive-shjd/conf:all-srcs",
        "//app/job/main/archive-shjd/http:all-srcs",
        "//app/job/main/archive-shjd/model:all-srcs",
        "//app/job/main/archive-shjd/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
