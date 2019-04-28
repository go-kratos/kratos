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
        "//app/service/main/coin/api:all-srcs",
        "//app/service/main/coin/cmd:all-srcs",
        "//app/service/main/coin/conf:all-srcs",
        "//app/service/main/coin/dao:all-srcs",
        "//app/service/main/coin/model:all-srcs",
        "//app/service/main/coin/server/gorpc:all-srcs",
        "//app/service/main/coin/server/grpc:all-srcs",
        "//app/service/main/coin/server/http:all-srcs",
        "//app/service/main/coin/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
