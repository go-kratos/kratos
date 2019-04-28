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
        "//app/admin/main/relation/cmd:all-srcs",
        "//app/admin/main/relation/conf:all-srcs",
        "//app/admin/main/relation/dao:all-srcs",
        "//app/admin/main/relation/http:all-srcs",
        "//app/admin/main/relation/model:all-srcs",
        "//app/admin/main/relation/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
