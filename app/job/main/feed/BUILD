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
        "//app/job/main/feed/cmd:all-srcs",
        "//app/job/main/feed/conf:all-srcs",
        "//app/job/main/feed/dao:all-srcs",
        "//app/job/main/feed/http:all-srcs",
        "//app/job/main/feed/model:all-srcs",
        "//app/job/main/feed/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
