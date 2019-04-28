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
        "//app/admin/main/videoup-task/cmd:all-srcs",
        "//app/admin/main/videoup-task/conf:all-srcs",
        "//app/admin/main/videoup-task/dao:all-srcs",
        "//app/admin/main/videoup-task/http:all-srcs",
        "//app/admin/main/videoup-task/model:all-srcs",
        "//app/admin/main/videoup-task/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
