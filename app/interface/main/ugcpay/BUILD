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
        "//app/interface/main/ugcpay/api/http:all-srcs",
        "//app/interface/main/ugcpay/cmd:all-srcs",
        "//app/interface/main/ugcpay/conf:all-srcs",
        "//app/interface/main/ugcpay/dao:all-srcs",
        "//app/interface/main/ugcpay/model:all-srcs",
        "//app/interface/main/ugcpay/server/http:all-srcs",
        "//app/interface/main/ugcpay/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
