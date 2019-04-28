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
        "//app/service/main/filter/api/grpc/v1:all-srcs",
        "//app/service/main/filter/cmd:all-srcs",
        "//app/service/main/filter/conf:all-srcs",
        "//app/service/main/filter/dao:all-srcs",
        "//app/service/main/filter/http:all-srcs",
        "//app/service/main/filter/model:all-srcs",
        "//app/service/main/filter/rpc/client:all-srcs",
        "//app/service/main/filter/rpc/server:all-srcs",
        "//app/service/main/filter/server/grpc:all-srcs",
        "//app/service/main/filter/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
