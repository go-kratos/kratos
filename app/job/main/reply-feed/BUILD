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
        "//app/job/main/reply-feed/cmd:all-srcs",
        "//app/job/main/reply-feed/conf:all-srcs",
        "//app/job/main/reply-feed/dao:all-srcs",
        "//app/job/main/reply-feed/model:all-srcs",
        "//app/job/main/reply-feed/server/http:all-srcs",
        "//app/job/main/reply-feed/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
