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
        "//app/admin/main/tv/cmd:all-srcs",
        "//app/admin/main/tv/conf:all-srcs",
        "//app/admin/main/tv/dao:all-srcs",
        "//app/admin/main/tv/http:all-srcs",
        "//app/admin/main/tv/model:all-srcs",
        "//app/admin/main/tv/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
