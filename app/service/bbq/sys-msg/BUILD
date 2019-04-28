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
        "//app/service/bbq/sys-msg/api/v1:all-srcs",
        "//app/service/bbq/sys-msg/cmd:all-srcs",
        "//app/service/bbq/sys-msg/internal/conf:all-srcs",
        "//app/service/bbq/sys-msg/internal/dao:all-srcs",
        "//app/service/bbq/sys-msg/internal/model:all-srcs",
        "//app/service/bbq/sys-msg/internal/server/grpc:all-srcs",
        "//app/service/bbq/sys-msg/internal/server/http:all-srcs",
        "//app/service/bbq/sys-msg/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
