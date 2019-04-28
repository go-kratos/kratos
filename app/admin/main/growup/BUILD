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
        "//app/admin/main/growup/cmd:all-srcs",
        "//app/admin/main/growup/conf:all-srcs",
        "//app/admin/main/growup/dao:all-srcs",
        "//app/admin/main/growup/http:all-srcs",
        "//app/admin/main/growup/model:all-srcs",
        "//app/admin/main/growup/service:all-srcs",
        "//app/admin/main/growup/util:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
