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
        "//app/service/live/gift/api/grpc/v1:all-srcs",
        "//app/service/live/gift/api/liverpc:all-srcs",
        "//app/service/live/gift/cmd:all-srcs",
        "//app/service/live/gift/conf:all-srcs",
        "//app/service/live/gift/dao:all-srcs",
        "//app/service/live/gift/model:all-srcs",
        "//app/service/live/gift/server/grpc:all-srcs",
        "//app/service/live/gift/server/http:all-srcs",
        "//app/service/live/gift/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
