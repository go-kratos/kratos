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
        "//app/service/bbq/notice-service/api/v1:all-srcs",
        "//app/service/bbq/notice-service/cmd:all-srcs",
        "//app/service/bbq/notice-service/internal/conf:all-srcs",
        "//app/service/bbq/notice-service/internal/dao:all-srcs",
        "//app/service/bbq/notice-service/internal/model:all-srcs",
        "//app/service/bbq/notice-service/internal/server/grpc:all-srcs",
        "//app/service/bbq/notice-service/internal/server/http:all-srcs",
        "//app/service/bbq/notice-service/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
