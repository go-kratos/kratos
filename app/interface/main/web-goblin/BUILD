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
        "//app/interface/main/web-goblin/cmd:all-srcs",
        "//app/interface/main/web-goblin/conf:all-srcs",
        "//app/interface/main/web-goblin/dao/share:all-srcs",
        "//app/interface/main/web-goblin/dao/web:all-srcs",
        "//app/interface/main/web-goblin/dao/wechat:all-srcs",
        "//app/interface/main/web-goblin/http:all-srcs",
        "//app/interface/main/web-goblin/model/share:all-srcs",
        "//app/interface/main/web-goblin/model/web:all-srcs",
        "//app/interface/main/web-goblin/model/wechat:all-srcs",
        "//app/interface/main/web-goblin/service/share:all-srcs",
        "//app/interface/main/web-goblin/service/web:all-srcs",
        "//app/interface/main/web-goblin/service/wechat:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
