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
        "//app/admin/main/dm/cmd:all-srcs",
        "//app/admin/main/dm/conf:all-srcs",
        "//app/admin/main/dm/dao:all-srcs",
        "//app/admin/main/dm/http:all-srcs",
        "//app/admin/main/dm/model:all-srcs",
        "//app/admin/main/dm/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
