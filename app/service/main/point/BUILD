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
        "//app/service/main/point/api:all-srcs",
        "//app/service/main/point/cmd:all-srcs",
        "//app/service/main/point/conf:all-srcs",
        "//app/service/main/point/dao:all-srcs",
        "//app/service/main/point/model:all-srcs",
        "//app/service/main/point/rpc/client:all-srcs",
        "//app/service/main/point/rpc/server:all-srcs",
        "//app/service/main/point/server/grpc:all-srcs",
        "//app/service/main/point/server/http:all-srcs",
        "//app/service/main/point/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
