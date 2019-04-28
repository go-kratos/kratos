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
        "//app/job/live/recommend-job/cmd:all-srcs",
        "//app/job/live/recommend-job/internal/conf:all-srcs",
        "//app/job/live/recommend-job/internal/dao:all-srcs",
        "//app/job/live/recommend-job/internal/model:all-srcs",
        "//app/job/live/recommend-job/internal/server/http:all-srcs",
        "//app/job/live/recommend-job/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
