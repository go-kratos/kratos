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
        "//app/service/bbq/recsys-recall/api/grpc/v1:all-srcs",
        "//app/service/bbq/recsys-recall/cmd:all-srcs",
        "//app/service/bbq/recsys-recall/conf:all-srcs",
        "//app/service/bbq/recsys-recall/dao:all-srcs",
        "//app/service/bbq/recsys-recall/model:all-srcs",
        "//app/service/bbq/recsys-recall/server/grpc:all-srcs",
        "//app/service/bbq/recsys-recall/server/http:all-srcs",
        "//app/service/bbq/recsys-recall/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
