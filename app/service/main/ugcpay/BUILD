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
        "//app/service/main/ugcpay/api/grpc/v1:all-srcs",
        "//app/service/main/ugcpay/api/http:all-srcs",
        "//app/service/main/ugcpay/cmd:all-srcs",
        "//app/service/main/ugcpay/conf:all-srcs",
        "//app/service/main/ugcpay/dao:all-srcs",
        "//app/service/main/ugcpay/model:all-srcs",
        "//app/service/main/ugcpay/server/grpc:all-srcs",
        "//app/service/main/ugcpay/server/http:all-srcs",
        "//app/service/main/ugcpay/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
