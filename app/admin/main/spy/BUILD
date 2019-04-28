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
        "//app/admin/main/spy/cmd:all-srcs",
        "//app/admin/main/spy/conf:all-srcs",
        "//app/admin/main/spy/dao:all-srcs",
        "//app/admin/main/spy/http:all-srcs",
        "//app/admin/main/spy/model:all-srcs",
        "//app/admin/main/spy/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
