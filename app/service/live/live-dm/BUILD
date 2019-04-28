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
        "//app/service/live/live-dm/api/grpc:all-srcs",
        "//app/service/live/live-dm/cmd:all-srcs",
        "//app/service/live/live-dm/conf:all-srcs",
        "//app/service/live/live-dm/dao:all-srcs",
        "//app/service/live/live-dm/model:all-srcs",
        "//app/service/live/live-dm/server/grpc:all-srcs",
        "//app/service/live/live-dm/server/http:all-srcs",
        "//app/service/live/live-dm/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
