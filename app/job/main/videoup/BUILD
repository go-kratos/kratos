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
        "//app/job/main/videoup/cmd:all-srcs",
        "//app/job/main/videoup/conf:all-srcs",
        "//app/job/main/videoup/dao:all-srcs",
        "//app/job/main/videoup/http:all-srcs",
        "//app/job/main/videoup/model/archive:all-srcs",
        "//app/job/main/videoup/model/manager:all-srcs",
        "//app/job/main/videoup/model/message:all-srcs",
        "//app/job/main/videoup/model/redis:all-srcs",
        "//app/job/main/videoup/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
