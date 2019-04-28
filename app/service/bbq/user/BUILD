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
        "//app/service/bbq/user/api:all-srcs",
        "//app/service/bbq/user/cmd:all-srcs",
        "//app/service/bbq/user/internal/conf:all-srcs",
        "//app/service/bbq/user/internal/dao:all-srcs",
        "//app/service/bbq/user/internal/model:all-srcs",
        "//app/service/bbq/user/internal/server/grpc:all-srcs",
        "//app/service/bbq/user/internal/server/http:all-srcs",
        "//app/service/bbq/user/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
