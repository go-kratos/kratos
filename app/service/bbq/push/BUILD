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
        "//app/service/bbq/push/api/grpc/v1:all-srcs",
        "//app/service/bbq/push/api/http/v1:all-srcs",
        "//app/service/bbq/push/cmd:all-srcs",
        "//app/service/bbq/push/conf:all-srcs",
        "//app/service/bbq/push/dao:all-srcs",
        "//app/service/bbq/push/model:all-srcs",
        "//app/service/bbq/push/server/grpc:all-srcs",
        "//app/service/bbq/push/server/http:all-srcs",
        "//app/service/bbq/push/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
