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
        "//app/admin/main/space/cmd:all-srcs",
        "//app/admin/main/space/conf:all-srcs",
        "//app/admin/main/space/dao:all-srcs",
        "//app/admin/main/space/http:all-srcs",
        "//app/admin/main/space/model:all-srcs",
        "//app/admin/main/space/service:all-srcs",
        "//app/admin/main/space/util:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
