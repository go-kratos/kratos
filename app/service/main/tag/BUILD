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
        "//app/service/main/tag/api:all-srcs",
        "//app/service/main/tag/cmd:all-srcs",
        "//app/service/main/tag/conf:all-srcs",
        "//app/service/main/tag/dao:all-srcs",
        "//app/service/main/tag/model:all-srcs",
        "//app/service/main/tag/rpc/client:all-srcs",
        "//app/service/main/tag/rpc/server:all-srcs",
        "//app/service/main/tag/server/grpc:all-srcs",
        "//app/service/main/tag/server/http:all-srcs",
        "//app/service/main/tag/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
