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
        "//app/job/main/block/cmd:all-srcs",
        "//app/job/main/block/conf:all-srcs",
        "//app/job/main/block/dao:all-srcs",
        "//app/job/main/block/http:all-srcs",
        "//app/job/main/block/model:all-srcs",
        "//app/job/main/block/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
