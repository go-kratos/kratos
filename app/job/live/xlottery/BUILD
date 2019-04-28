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
        "//app/job/live/xlottery/cmd:all-srcs",
        "//app/job/live/xlottery/internal/conf:all-srcs",
        "//app/job/live/xlottery/internal/dao:all-srcs",
        "//app/job/live/xlottery/internal/model:all-srcs",
        "//app/job/live/xlottery/internal/server/http:all-srcs",
        "//app/job/live/xlottery/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
