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
        "//app/interface/main/ugcpay-rank/api/http:all-srcs",
        "//app/interface/main/ugcpay-rank/cmd:all-srcs",
        "//app/interface/main/ugcpay-rank/internal/conf:all-srcs",
        "//app/interface/main/ugcpay-rank/internal/dao:all-srcs",
        "//app/interface/main/ugcpay-rank/internal/server/http:all-srcs",
        "//app/interface/main/ugcpay-rank/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
