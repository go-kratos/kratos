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
        "//app/service/main/push/api/gorpc:all-srcs",
        "//app/service/main/push/api/grpc/v1:all-srcs",
        "//app/service/main/push/cmd:all-srcs",
        "//app/service/main/push/conf:all-srcs",
        "//app/service/main/push/dao:all-srcs",
        "//app/service/main/push/model:all-srcs",
        "//app/service/main/push/server/gorpc:all-srcs",
        "//app/service/main/push/server/grpc:all-srcs",
        "//app/service/main/push/server/http:all-srcs",
        "//app/service/main/push/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
