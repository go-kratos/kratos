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
        "//app/service/live/xrewardcenter/api/grpc/v1:all-srcs",
        "//app/service/live/xrewardcenter/cmd:all-srcs",
        "//app/service/live/xrewardcenter/conf:all-srcs",
        "//app/service/live/xrewardcenter/dao:all-srcs",
        "//app/service/live/xrewardcenter/model:all-srcs",
        "//app/service/live/xrewardcenter/server/grpc:all-srcs",
        "//app/service/live/xrewardcenter/server/http:all-srcs",
        "//app/service/live/xrewardcenter/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
