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
        "//app/service/live/xcaptcha/api/grpc/v1:all-srcs",
        "//app/service/live/xcaptcha/api/http/v1:all-srcs",
        "//app/service/live/xcaptcha/cmd:all-srcs",
        "//app/service/live/xcaptcha/conf:all-srcs",
        "//app/service/live/xcaptcha/dao:all-srcs",
        "//app/service/live/xcaptcha/model:all-srcs",
        "//app/service/live/xcaptcha/server/grpc:all-srcs",
        "//app/service/live/xcaptcha/server/http:all-srcs",
        "//app/service/live/xcaptcha/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
