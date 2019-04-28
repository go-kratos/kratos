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
        "//app/interface/openplatform/article:all-srcs",
        "//app/interface/openplatform/monitor-end/cmd:all-srcs",
        "//app/interface/openplatform/monitor-end/conf:all-srcs",
        "//app/interface/openplatform/monitor-end/dao:all-srcs",
        "//app/interface/openplatform/monitor-end/http:all-srcs",
        "//app/interface/openplatform/monitor-end/model:all-srcs",
        "//app/interface/openplatform/monitor-end/service:all-srcs",
        "//app/interface/openplatform/seo:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
