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
        "//app/service/main/thumbup/api:all-srcs",
        "//app/service/main/thumbup/cmd:all-srcs",
        "//app/service/main/thumbup/conf:all-srcs",
        "//app/service/main/thumbup/dao:all-srcs",
        "//app/service/main/thumbup/model:all-srcs",
        "//app/service/main/thumbup/rpc/client:all-srcs",
        "//app/service/main/thumbup/server/gorpc:all-srcs",
        "//app/service/main/thumbup/server/grpc:all-srcs",
        "//app/service/main/thumbup/server/http:all-srcs",
        "//app/service/main/thumbup/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
