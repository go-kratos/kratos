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
        "//app/service/main/rank/api/gorpc:all-srcs",
        "//app/service/main/rank/cmd:all-srcs",
        "//app/service/main/rank/conf:all-srcs",
        "//app/service/main/rank/dao:all-srcs",
        "//app/service/main/rank/model:all-srcs",
        "//app/service/main/rank/server/gorpc:all-srcs",
        "//app/service/main/rank/server/grpc:all-srcs",
        "//app/service/main/rank/server/http:all-srcs",
        "//app/service/main/rank/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
