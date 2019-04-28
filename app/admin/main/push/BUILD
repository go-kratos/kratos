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
        "//app/admin/main/push/cmd:all-srcs",
        "//app/admin/main/push/conf:all-srcs",
        "//app/admin/main/push/dao:all-srcs",
        "//app/admin/main/push/http:all-srcs",
        "//app/admin/main/push/model:all-srcs",
        "//app/admin/main/push/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
