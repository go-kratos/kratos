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
        "//app/admin/main/macross/cmd:all-srcs",
        "//app/admin/main/macross/conf:all-srcs",
        "//app/admin/main/macross/dao:all-srcs",
        "//app/admin/main/macross/http:all-srcs",
        "//app/admin/main/macross/model/errors:all-srcs",
        "//app/admin/main/macross/model/mail:all-srcs",
        "//app/admin/main/macross/model/manager:all-srcs",
        "//app/admin/main/macross/model/package:all-srcs",
        "//app/admin/main/macross/model/publish:all-srcs",
        "//app/admin/main/macross/service:all-srcs",
        "//app/admin/main/macross/tools:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
