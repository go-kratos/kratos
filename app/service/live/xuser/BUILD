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
        "//app/service/live/xuser/api/grpc:all-srcs",
        "//app/service/live/xuser/cmd:all-srcs",
        "//app/service/live/xuser/conf:all-srcs",
        "//app/service/live/xuser/dao:all-srcs",
        "//app/service/live/xuser/model:all-srcs",
        "//app/service/live/xuser/server/grpc:all-srcs",
        "//app/service/live/xuser/server/http:all-srcs",
        "//app/service/live/xuser/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
