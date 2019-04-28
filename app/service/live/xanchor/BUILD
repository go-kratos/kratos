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
        "//app/service/live/xanchor/api/grpc/v1:all-srcs",
        "//app/service/live/xanchor/cmd:all-srcs",
        "//app/service/live/xanchor/conf:all-srcs",
        "//app/service/live/xanchor/dao:all-srcs",
        "//app/service/live/xanchor/model:all-srcs",
        "//app/service/live/xanchor/server/grpc:all-srcs",
        "//app/service/live/xanchor/server/http:all-srcs",
        "//app/service/live/xanchor/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
