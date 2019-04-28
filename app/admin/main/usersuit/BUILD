package(default_visibility = ["//visibility:public"])

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
        "//app/admin/main/usersuit/cmd:all-srcs",
        "//app/admin/main/usersuit/conf:all-srcs",
        "//app/admin/main/usersuit/dao:all-srcs",
        "//app/admin/main/usersuit/http:all-srcs",
        "//app/admin/main/usersuit/model:all-srcs",
        "//app/admin/main/usersuit/service:all-srcs",
    ],
    tags = ["automanaged"],
)
