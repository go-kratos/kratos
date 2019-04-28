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
        "//app/job/bbq/cms/cmd:all-srcs",
        "//app/job/bbq/cms/internal/dao:all-srcs",
        "//app/job/bbq/cms/internal/model:all-srcs",
        "//app/job/bbq/cms/internal/server/http:all-srcs",
        "//app/job/bbq/cms/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
