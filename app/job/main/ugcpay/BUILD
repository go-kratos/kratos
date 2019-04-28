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
        "//app/job/main/ugcpay/cmd:all-srcs",
        "//app/job/main/ugcpay/conf:all-srcs",
        "//app/job/main/ugcpay/dao:all-srcs",
        "//app/job/main/ugcpay/model:all-srcs",
        "//app/job/main/ugcpay/server/http:all-srcs",
        "//app/job/main/ugcpay/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
