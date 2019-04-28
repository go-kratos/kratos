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
        "//app/interface/main/app-feed/cmd:all-srcs",
        "//app/interface/main/app-feed/conf:all-srcs",
        "//app/interface/main/app-feed/dao:all-srcs",
        "//app/interface/main/app-feed/http:all-srcs",
        "//app/interface/main/app-feed/model:all-srcs",
        "//app/interface/main/app-feed/service/external:all-srcs",
        "//app/interface/main/app-feed/service/feed:all-srcs",
        "//app/interface/main/app-feed/service/ping:all-srcs",
        "//app/interface/main/app-feed/service/region:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
