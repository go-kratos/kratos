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
        "//app/job/main/tag/cmd:all-srcs",
        "//app/job/main/tag/conf:all-srcs",
        "//app/job/main/tag/dao:all-srcs",
        "//app/job/main/tag/http:all-srcs",
        "//app/job/main/tag/model:all-srcs",
        "//app/job/main/tag/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
