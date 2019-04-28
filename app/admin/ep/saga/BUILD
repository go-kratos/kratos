package(default_visibility = ["//visibility:public"])

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
        "//app/admin/ep/saga/api/grpc/v1:all-srcs",
        "//app/admin/ep/saga/cmd:all-srcs",
        "//app/admin/ep/saga/conf:all-srcs",
        "//app/admin/ep/saga/dao:all-srcs",
        "//app/admin/ep/saga/http:all-srcs",
        "//app/admin/ep/saga/model:all-srcs",
        "//app/admin/ep/saga/server/grpc:all-srcs",
        "//app/admin/ep/saga/service:all-srcs",
    ],
    tags = ["automanaged"],
)
