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
        "//app/service/main/spy/api:all-srcs",
        "//app/service/main/spy/cmd:all-srcs",
        "//app/service/main/spy/conf:all-srcs",
        "//app/service/main/spy/dao:all-srcs",
        "//app/service/main/spy/model:all-srcs",
        "//app/service/main/spy/rpc/client:all-srcs",
        "//app/service/main/spy/rpc/server:all-srcs",
        "//app/service/main/spy/server/grpc:all-srcs",
        "//app/service/main/spy/server/http:all-srcs",
        "//app/service/main/spy/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
