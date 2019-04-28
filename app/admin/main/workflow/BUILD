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
        "//app/admin/main/workflow/cmd:all-srcs",
        "//app/admin/main/workflow/dao:all-srcs",
        "//app/admin/main/workflow/http:all-srcs",
        "//app/admin/main/workflow/model:all-srcs",
        "//app/admin/main/workflow/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
