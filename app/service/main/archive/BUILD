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
        "//app/service/main/archive/api:all-srcs",
        "//app/service/main/archive/cmd:all-srcs",
        "//app/service/main/archive/conf:all-srcs",
        "//app/service/main/archive/dao:all-srcs",
        "//app/service/main/archive/model:all-srcs",
        "//app/service/main/archive/server/gorpc:all-srcs",
        "//app/service/main/archive/server/grpc:all-srcs",
        "//app/service/main/archive/server/http:all-srcs",
        "//app/service/main/archive/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
