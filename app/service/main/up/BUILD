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
        "//app/service/main/up/api/gorpc:all-srcs",
        "//app/service/main/up/api/v1:all-srcs",
        "//app/service/main/up/cmd:all-srcs",
        "//app/service/main/up/conf:all-srcs",
        "//app/service/main/up/dao:all-srcs",
        "//app/service/main/up/model:all-srcs",
        "//app/service/main/up/server/gorpc:all-srcs",
        "//app/service/main/up/server/grpc:all-srcs",
        "//app/service/main/up/server/http:all-srcs",
        "//app/service/main/up/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
