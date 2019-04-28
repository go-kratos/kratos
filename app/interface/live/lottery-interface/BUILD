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
        "//app/interface/live/lottery-interface/cmd:all-srcs",
        "//app/interface/live/lottery-interface/internal/conf:all-srcs",
        "//app/interface/live/lottery-interface/internal/dao:all-srcs",
        "//app/interface/live/lottery-interface/internal/model:all-srcs",
        "//app/interface/live/lottery-interface/internal/server/http:all-srcs",
        "//app/interface/live/lottery-interface/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
