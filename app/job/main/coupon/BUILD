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
        "//app/job/main/coupon/cmd:all-srcs",
        "//app/job/main/coupon/conf:all-srcs",
        "//app/job/main/coupon/dao:all-srcs",
        "//app/job/main/coupon/http:all-srcs",
        "//app/job/main/coupon/model:all-srcs",
        "//app/job/main/coupon/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
