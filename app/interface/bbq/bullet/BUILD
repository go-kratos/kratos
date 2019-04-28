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
        "//app/interface/bbq/bullet/api:all-srcs",
        "//app/interface/bbq/bullet/cmd:all-srcs",
        "//app/interface/bbq/bullet/internal/conf:all-srcs",
        "//app/interface/bbq/bullet/internal/dao:all-srcs",
        "//app/interface/bbq/bullet/internal/model:all-srcs",
        "//app/interface/bbq/bullet/internal/server/http:all-srcs",
        "//app/interface/bbq/bullet/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
