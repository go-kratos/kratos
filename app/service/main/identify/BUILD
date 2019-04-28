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
        "//app/service/main/identify/api/grpc:all-srcs",
        "//app/service/main/identify/cmd:all-srcs",
        "//app/service/main/identify/conf:all-srcs",
        "//app/service/main/identify/dao:all-srcs",
        "//app/service/main/identify/model:all-srcs",
        "//app/service/main/identify/server/grpc:all-srcs",
        "//app/service/main/identify/server/http:all-srcs",
        "//app/service/main/identify/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
