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
        "//app/service/main/figure/cmd:all-srcs",
        "//app/service/main/figure/conf:all-srcs",
        "//app/service/main/figure/dao:all-srcs",
        "//app/service/main/figure/http:all-srcs",
        "//app/service/main/figure/model:all-srcs",
        "//app/service/main/figure/rpc/client:all-srcs",
        "//app/service/main/figure/rpc/server:all-srcs",
        "//app/service/main/figure/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
