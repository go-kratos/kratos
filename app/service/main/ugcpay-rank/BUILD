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
        "//app/service/main/ugcpay-rank/api/http:all-srcs",
        "//app/service/main/ugcpay-rank/api/v1:all-srcs",
        "//app/service/main/ugcpay-rank/cmd:all-srcs",
        "//app/service/main/ugcpay-rank/internal/conf:all-srcs",
        "//app/service/main/ugcpay-rank/internal/dao:all-srcs",
        "//app/service/main/ugcpay-rank/internal/model:all-srcs",
        "//app/service/main/ugcpay-rank/internal/server/grpc:all-srcs",
        "//app/service/main/ugcpay-rank/internal/server/http:all-srcs",
        "//app/service/main/ugcpay-rank/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
