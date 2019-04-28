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
        "//app/service/main/account/api:all-srcs",
        "//app/service/main/account/cmd:all-srcs",
        "//app/service/main/account/conf:all-srcs",
        "//app/service/main/account/dao:all-srcs",
        "//app/service/main/account/model:all-srcs",
        "//app/service/main/account/rpc:all-srcs",
        "//app/service/main/account/server/grpc:all-srcs",
        "//app/service/main/account/server/http:all-srcs",
        "//app/service/main/account/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
