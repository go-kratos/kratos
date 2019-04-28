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
        "//app/service/bbq/topic/api:all-srcs",
        "//app/service/bbq/topic/cmd:all-srcs",
        "//app/service/bbq/topic/internal/dao:all-srcs",
        "//app/service/bbq/topic/internal/model:all-srcs",
        "//app/service/bbq/topic/internal/server/grpc:all-srcs",
        "//app/service/bbq/topic/internal/server/http:all-srcs",
        "//app/service/bbq/topic/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
