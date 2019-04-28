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
        "//app/service/main/coupon/api:all-srcs",
        "//app/service/main/coupon/cmd:all-srcs",
        "//app/service/main/coupon/conf:all-srcs",
        "//app/service/main/coupon/dao:all-srcs",
        "//app/service/main/coupon/http:all-srcs",
        "//app/service/main/coupon/model:all-srcs",
        "//app/service/main/coupon/rpc/client:all-srcs",
        "//app/service/main/coupon/rpc/server:all-srcs",
        "//app/service/main/coupon/server/grpc:all-srcs",
        "//app/service/main/coupon/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
