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
        "//app/admin/main/up/cmd:all-srcs",
        "//app/admin/main/up/conf:all-srcs",
        "//app/admin/main/up/dao:all-srcs",
        "//app/admin/main/up/http:all-srcs",
        "//app/admin/main/up/model:all-srcs",
        "//app/admin/main/up/service:all-srcs",
        "//app/admin/main/up/util:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
