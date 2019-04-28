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
        "//app/admin/main/manager/api:all-srcs",
        "//app/admin/main/manager/cmd:all-srcs",
        "//app/admin/main/manager/conf:all-srcs",
        "//app/admin/main/manager/dao:all-srcs",
        "//app/admin/main/manager/model:all-srcs",
        "//app/admin/main/manager/server/grpc:all-srcs",
        "//app/admin/main/manager/server/http:all-srcs",
        "//app/admin/main/manager/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
