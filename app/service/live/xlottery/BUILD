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
        "//app/service/live/xlottery/api/grpc/v1:all-srcs",
        "//app/service/live/xlottery/cmd:all-srcs",
        "//app/service/live/xlottery/conf:all-srcs",
        "//app/service/live/xlottery/dao:all-srcs",
        "//app/service/live/xlottery/model:all-srcs",
        "//app/service/live/xlottery/server/grpc:all-srcs",
        "//app/service/live/xlottery/server/http:all-srcs",
        "//app/service/live/xlottery/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
