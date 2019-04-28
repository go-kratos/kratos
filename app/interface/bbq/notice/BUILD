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
        "//app/interface/bbq/notice/cmd:all-srcs",
        "//app/interface/bbq/notice/internal/dao:all-srcs",
        "//app/interface/bbq/notice/internal/model:all-srcs",
        "//app/interface/bbq/notice/internal/server/http:all-srcs",
        "//app/interface/bbq/notice/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
