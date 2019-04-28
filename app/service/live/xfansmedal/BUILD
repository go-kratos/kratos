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
        "//app/service/live/xfansmedal/api/v1:all-srcs",
        "//app/service/live/xfansmedal/cmd:all-srcs",
        "//app/service/live/xfansmedal/internal/conf:all-srcs",
        "//app/service/live/xfansmedal/internal/dao:all-srcs",
        "//app/service/live/xfansmedal/internal/model:all-srcs",
        "//app/service/live/xfansmedal/internal/server/grpc:all-srcs",
        "//app/service/live/xfansmedal/internal/server/http:all-srcs",
        "//app/service/live/xfansmedal/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
