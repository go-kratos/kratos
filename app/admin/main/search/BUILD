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
        "//app/admin/main/search/cmd:all-srcs",
        "//app/admin/main/search/conf:all-srcs",
        "//app/admin/main/search/dao:all-srcs",
        "//app/admin/main/search/http:all-srcs",
        "//app/admin/main/search/model:all-srcs",
        "//app/admin/main/search/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
