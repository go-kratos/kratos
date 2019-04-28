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
        "//app/service/bbq/recsys/api/grpc/v1:all-srcs",
        "//app/service/bbq/recsys/cmd:all-srcs",
        "//app/service/bbq/recsys/conf:all-srcs",
        "//app/service/bbq/recsys/dao:all-srcs",
        "//app/service/bbq/recsys/model:all-srcs",
        "//app/service/bbq/recsys/server/grpc:all-srcs",
        "//app/service/bbq/recsys/server/http:all-srcs",
        "//app/service/bbq/recsys/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
