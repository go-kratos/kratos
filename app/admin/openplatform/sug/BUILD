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
        "//app/admin/openplatform/sug/cmd:all-srcs",
        "//app/admin/openplatform/sug/conf:all-srcs",
        "//app/admin/openplatform/sug/dao:all-srcs",
        "//app/admin/openplatform/sug/http:all-srcs",
        "//app/admin/openplatform/sug/model:all-srcs",
        "//app/admin/openplatform/sug/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
