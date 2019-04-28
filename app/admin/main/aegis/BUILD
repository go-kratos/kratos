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
        "//app/admin/main/aegis/cmd:all-srcs",
        "//app/admin/main/aegis/conf:all-srcs",
        "//app/admin/main/aegis/dao:all-srcs",
        "//app/admin/main/aegis/model:all-srcs",
        "//app/admin/main/aegis/server/databus:all-srcs",
        "//app/admin/main/aegis/server/http:all-srcs",
        "//app/admin/main/aegis/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
