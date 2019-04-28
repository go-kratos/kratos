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
        "//app/service/main/feed/cmd:all-srcs",
        "//app/service/main/feed/conf:all-srcs",
        "//app/service/main/feed/dao:all-srcs",
        "//app/service/main/feed/http:all-srcs",
        "//app/service/main/feed/model:all-srcs",
        "//app/service/main/feed/rpc/client:all-srcs",
        "//app/service/main/feed/rpc/server:all-srcs",
        "//app/service/main/feed/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
