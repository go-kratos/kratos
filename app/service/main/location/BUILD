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
        "//app/service/main/location/api:all-srcs",
        "//app/service/main/location/cmd:all-srcs",
        "//app/service/main/location/conf:all-srcs",
        "//app/service/main/location/dao:all-srcs",
        "//app/service/main/location/model:all-srcs",
        "//app/service/main/location/rpc/client:all-srcs",
        "//app/service/main/location/rpc/server:all-srcs",
        "//app/service/main/location/server/grpc:all-srcs",
        "//app/service/main/location/server/http:all-srcs",
        "//app/service/main/location/service:all-srcs",
        "//app/service/main/location/zlimit:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
