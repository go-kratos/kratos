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
        "//app/job/bbq/recall/api/v1:all-srcs",
        "//app/job/bbq/recall/cmd:all-srcs",
        "//app/job/bbq/recall/internal/conf:all-srcs",
        "//app/job/bbq/recall/internal/dao:all-srcs",
        "//app/job/bbq/recall/internal/model:all-srcs",
        "//app/job/bbq/recall/internal/service:all-srcs",
        "//app/job/bbq/recall/proto:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
