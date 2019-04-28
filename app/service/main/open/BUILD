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
        "//app/service/main/open/api/grpc/v1:all-srcs",
        "//app/service/main/open/cmd:all-srcs",
        "//app/service/main/open/conf:all-srcs",
        "//app/service/main/open/dao:all-srcs",
        "//app/service/main/open/model:all-srcs",
        "//app/service/main/open/server/grpc:all-srcs",
        "//app/service/main/open/server/http:all-srcs",
        "//app/service/main/open/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
