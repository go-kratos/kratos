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
        "//app/service/bbq/video/api/grpc/v1:all-srcs",
        "//app/service/bbq/video/api/http/v1:all-srcs",
        "//app/service/bbq/video/cmd:all-srcs",
        "//app/service/bbq/video/conf:all-srcs",
        "//app/service/bbq/video/dao:all-srcs",
        "//app/service/bbq/video/model:all-srcs",
        "//app/service/bbq/video/server/grpc:all-srcs",
        "//app/service/bbq/video/server/http:all-srcs",
        "//app/service/bbq/video/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
