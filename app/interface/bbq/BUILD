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
        "//app/interface/bbq/app-bbq:all-srcs",
        "//app/interface/bbq/bullet:all-srcs",
        "//app/interface/bbq/comment:all-srcs",
        "//app/interface/bbq/common:all-srcs",
        "//app/interface/bbq/feed:all-srcs",
        "//app/interface/bbq/notice:all-srcs",
        "//app/interface/bbq/search:all-srcs",
        "//app/interface/bbq/topic:all-srcs",
        "//app/interface/bbq/user:all-srcs",
        "//app/interface/bbq/video:all-srcs",
        "//app/interface/bbq/wechat:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
