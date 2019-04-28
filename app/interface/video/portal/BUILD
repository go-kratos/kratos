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
        "//app/interface/video/portal/cmd:all-srcs",
        "//app/interface/video/portal/conf:all-srcs",
        "//app/interface/video/portal/dao:all-srcs",
        "//app/interface/video/portal/middleware:all-srcs",
        "//app/interface/video/portal/model:all-srcs",
        "//app/interface/video/portal/server/http:all-srcs",
        "//app/interface/video/portal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
