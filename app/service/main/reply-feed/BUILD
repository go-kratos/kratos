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
        "//app/service/main/reply-feed/api:all-srcs",
        "//app/service/main/reply-feed/cmd:all-srcs",
        "//app/service/main/reply-feed/conf:all-srcs",
        "//app/service/main/reply-feed/dao:all-srcs",
        "//app/service/main/reply-feed/model:all-srcs",
        "//app/service/main/reply-feed/server/grpc:all-srcs",
        "//app/service/main/reply-feed/server/http:all-srcs",
        "//app/service/main/reply-feed/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
