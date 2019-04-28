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
        "//app/service/bbq/comment/api:all-srcs",
        "//app/service/bbq/comment/cmd:all-srcs",
        "//app/service/bbq/comment/internal/dao:all-srcs",
        "//app/service/bbq/comment/internal/model:all-srcs",
        "//app/service/bbq/comment/internal/server/grpc:all-srcs",
        "//app/service/bbq/comment/internal/server/http:all-srcs",
        "//app/service/bbq/comment/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
