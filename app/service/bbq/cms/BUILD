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
        "//app/service/bbq/cms/cmd:all-srcs",
        "//app/service/bbq/cms/internal/dao:all-srcs",
        "//app/service/bbq/cms/internal/model:all-srcs",
        "//app/service/bbq/cms/internal/server/http:all-srcs",
        "//app/service/bbq/cms/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
