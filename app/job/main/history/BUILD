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
        "//app/job/main/history/cmd:all-srcs",
        "//app/job/main/history/conf:all-srcs",
        "//app/job/main/history/dao:all-srcs",
        "//app/job/main/history/http:all-srcs",
        "//app/job/main/history/model:all-srcs",
        "//app/job/main/history/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
