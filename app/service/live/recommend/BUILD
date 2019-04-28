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
        "//app/service/live/recommend/api/grpc/v1:all-srcs",
        "//app/service/live/recommend/cmd:all-srcs",
        "//app/service/live/recommend/internal/conf:all-srcs",
        "//app/service/live/recommend/internal/dao:all-srcs",
        "//app/service/live/recommend/internal/model:all-srcs",
        "//app/service/live/recommend/internal/server/grpc:all-srcs",
        "//app/service/live/recommend/internal/server/http:all-srcs",
        "//app/service/live/recommend/internal/service:all-srcs",
        "//app/service/live/recommend/recconst:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
