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
        "//app/service/bbq/cms:all-srcs",
        "//app/service/bbq/comment:all-srcs",
        "//app/service/bbq/common:all-srcs",
        "//app/service/bbq/notice-service:all-srcs",
        "//app/service/bbq/push:all-srcs",
        "//app/service/bbq/recsys:all-srcs",
        "//app/service/bbq/recsys-recall:all-srcs",
        "//app/service/bbq/search:all-srcs",
        "//app/service/bbq/sys-msg:all-srcs",
        "//app/service/bbq/topic:all-srcs",
        "//app/service/bbq/user:all-srcs",
        "//app/service/bbq/video:all-srcs",
        "//app/service/bbq/video-image:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
