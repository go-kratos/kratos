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
        "//app/interface/bbq/wechat/cmd:all-srcs",
        "//app/interface/bbq/wechat/dao:all-srcs",
        "//app/interface/bbq/wechat/internal/conf:all-srcs",
        "//app/interface/bbq/wechat/internal/model:all-srcs",
        "//app/interface/bbq/wechat/internal/server/http:all-srcs",
        "//app/interface/bbq/wechat/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
