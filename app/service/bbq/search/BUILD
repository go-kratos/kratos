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
        "//app/service/bbq/search/api/grpc/v1:all-srcs",
        "//app/service/bbq/search/cmd:all-srcs",
        "//app/service/bbq/search/conf:all-srcs",
        "//app/service/bbq/search/dao:all-srcs",
        "//app/service/bbq/search/model:all-srcs",
        "//app/service/bbq/search/server/grpc:all-srcs",
        "//app/service/bbq/search/server/http:all-srcs",
        "//app/service/bbq/search/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
