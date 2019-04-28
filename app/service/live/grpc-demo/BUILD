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
        "//app/service/live/grpc-demo/api/grpc/v1:all-srcs",
        "//app/service/live/grpc-demo/cmd:all-srcs",
        "//app/service/live/grpc-demo/conf:all-srcs",
        "//app/service/live/grpc-demo/dao:all-srcs",
        "//app/service/live/grpc-demo/model:all-srcs",
        "//app/service/live/grpc-demo/server/grpc:all-srcs",
        "//app/service/live/grpc-demo/server/http:all-srcs",
        "//app/service/live/grpc-demo/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
