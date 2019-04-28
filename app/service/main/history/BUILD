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
        "//app/service/main/history/api/grpc:all-srcs",
        "//app/service/main/history/cmd:all-srcs",
        "//app/service/main/history/conf:all-srcs",
        "//app/service/main/history/dao:all-srcs",
        "//app/service/main/history/model:all-srcs",
        "//app/service/main/history/server/grpc:all-srcs",
        "//app/service/main/history/server/http:all-srcs",
        "//app/service/main/history/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
