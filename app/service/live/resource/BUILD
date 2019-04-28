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
        "//app/service/live/resource/api/grpc/v1:all-srcs",
        "//app/service/live/resource/api/grpc/v2:all-srcs",
        "//app/service/live/resource/api/http/v1:all-srcs",
        "//app/service/live/resource/cmd:all-srcs",
        "//app/service/live/resource/conf:all-srcs",
        "//app/service/live/resource/dao:all-srcs",
        "//app/service/live/resource/lrucache:all-srcs",
        "//app/service/live/resource/model:all-srcs",
        "//app/service/live/resource/sdk:all-srcs",
        "//app/service/live/resource/server/grpc:all-srcs",
        "//app/service/live/resource/server/http:all-srcs",
        "//app/service/live/resource/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
