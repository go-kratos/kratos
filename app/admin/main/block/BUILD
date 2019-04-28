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
        "//app/admin/main/block/cmd:all-srcs",
        "//app/admin/main/block/conf:all-srcs",
        "//app/admin/main/block/dao:all-srcs",
        "//app/admin/main/block/http:all-srcs",
        "//app/admin/main/block/model:all-srcs",
        "//app/admin/main/block/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
