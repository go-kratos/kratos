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
        "//app/interface/live/app-ucenter/api/http/v1:all-srcs",
        "//app/interface/live/app-ucenter/cmd:all-srcs",
        "//app/interface/live/app-ucenter/conf:all-srcs",
        "//app/interface/live/app-ucenter/dao:all-srcs",
        "//app/interface/live/app-ucenter/model:all-srcs",
        "//app/interface/live/app-ucenter/server/http:all-srcs",
        "//app/interface/live/app-ucenter/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
