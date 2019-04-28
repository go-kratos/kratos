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
        "//app/interface/live/web-ucenter/api/http:all-srcs",
        "//app/interface/live/web-ucenter/cmd:all-srcs",
        "//app/interface/live/web-ucenter/conf:all-srcs",
        "//app/interface/live/web-ucenter/dao:all-srcs",
        "//app/interface/live/web-ucenter/model:all-srcs",
        "//app/interface/live/web-ucenter/server/http:all-srcs",
        "//app/interface/live/web-ucenter/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
