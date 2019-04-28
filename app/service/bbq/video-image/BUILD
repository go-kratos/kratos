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
        "//app/service/bbq/video-image/api/grpc/v1:all-srcs",
        "//app/service/bbq/video-image/api/http/v1:all-srcs",
        "//app/service/bbq/video-image/cmd:all-srcs",
        "//app/service/bbq/video-image/conf:all-srcs",
        "//app/service/bbq/video-image/dao:all-srcs",
        "//app/service/bbq/video-image/model:all-srcs",
        "//app/service/bbq/video-image/server/grpc:all-srcs",
        "//app/service/bbq/video-image/server/http:all-srcs",
        "//app/service/bbq/video-image/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
