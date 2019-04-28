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
        "//app/service/main/vipinfo/api:all-srcs",
        "//app/service/main/vipinfo/cmd:all-srcs",
        "//app/service/main/vipinfo/conf:all-srcs",
        "//app/service/main/vipinfo/dao:all-srcs",
        "//app/service/main/vipinfo/model:all-srcs",
        "//app/service/main/vipinfo/server/grpc:all-srcs",
        "//app/service/main/vipinfo/server/http:all-srcs",
        "//app/service/main/vipinfo/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
