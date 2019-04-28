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
        "//app/service/main/vip/api:all-srcs",
        "//app/service/main/vip/cmd:all-srcs",
        "//app/service/main/vip/conf:all-srcs",
        "//app/service/main/vip/dao:all-srcs",
        "//app/service/main/vip/http:all-srcs",
        "//app/service/main/vip/model:all-srcs",
        "//app/service/main/vip/rpc/client:all-srcs",
        "//app/service/main/vip/rpc/server:all-srcs",
        "//app/service/main/vip/server/grpc:all-srcs",
        "//app/service/main/vip/service:all-srcs",
        "//app/service/main/vip/verify:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
