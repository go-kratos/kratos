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
        "//app/admin/main/open/cmd:all-srcs",
        "//app/admin/main/open/conf:all-srcs",
        "//app/admin/main/open/dao:all-srcs",
        "//app/admin/main/open/http:all-srcs",
        "//app/admin/main/open/model:all-srcs",
        "//app/admin/main/open/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
