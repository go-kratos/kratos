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
        "//app/job/main/mcn/cmd:all-srcs",
        "//app/job/main/mcn/conf:all-srcs",
        "//app/job/main/mcn/dao:all-srcs",
        "//app/job/main/mcn/model:all-srcs",
        "//app/job/main/mcn/server/http:all-srcs",
        "//app/job/main/mcn/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
