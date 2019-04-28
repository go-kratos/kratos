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
        "//app/job/main/point/cmd:all-srcs",
        "//app/job/main/point/conf:all-srcs",
        "//app/job/main/point/dao:all-srcs",
        "//app/job/main/point/http:all-srcs",
        "//app/job/main/point/model:all-srcs",
        "//app/job/main/point/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
