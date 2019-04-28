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
        "//app/service/main/resource/api/v1:all-srcs",
        "//app/service/main/resource/cmd:all-srcs",
        "//app/service/main/resource/conf:all-srcs",
        "//app/service/main/resource/dao/abtest:all-srcs",
        "//app/service/main/resource/dao/ads:all-srcs",
        "//app/service/main/resource/dao/alarm:all-srcs",
        "//app/service/main/resource/dao/cpm:all-srcs",
        "//app/service/main/resource/dao/manager:all-srcs",
        "//app/service/main/resource/dao/resource:all-srcs",
        "//app/service/main/resource/dao/show:all-srcs",
        "//app/service/main/resource/http:all-srcs",
        "//app/service/main/resource/model:all-srcs",
        "//app/service/main/resource/rpc/client:all-srcs",
        "//app/service/main/resource/rpc/server:all-srcs",
        "//app/service/main/resource/server/grpc:all-srcs",
        "//app/service/main/resource/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
