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
        "//app/admin/main/answer/cmd:all-srcs",
        "//app/admin/main/answer/conf:all-srcs",
        "//app/admin/main/answer/dao:all-srcs",
        "//app/admin/main/answer/http:all-srcs",
        "//app/admin/main/answer/model:all-srcs",
        "//app/admin/main/answer/service:all-srcs",
    ],
    tags = ["automanaged"],
)
