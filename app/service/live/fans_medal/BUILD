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
        "//app/service/live/fans_medal/api/liverpc:all-srcs",
        "//app/service/live/fans_medal/cmd:all-srcs",
        "//app/service/live/fans_medal/conf:all-srcs",
        "//app/service/live/fans_medal/dao:all-srcs",
        "//app/service/live/fans_medal/model:all-srcs",
        "//app/service/live/fans_medal/server/http:all-srcs",
        "//app/service/live/fans_medal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
