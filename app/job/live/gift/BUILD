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
        "//app/job/live/gift/api:all-srcs",
        "//app/job/live/gift/cmd:all-srcs",
        "//app/job/live/gift/internal/conf:all-srcs",
        "//app/job/live/gift/internal/dao:all-srcs",
        "//app/job/live/gift/internal/model:all-srcs",
        "//app/job/live/gift/internal/server/http:all-srcs",
        "//app/job/live/gift/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
