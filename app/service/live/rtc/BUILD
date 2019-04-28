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
        "//app/service/live/rtc/api/v1:all-srcs",
        "//app/service/live/rtc/cmd:all-srcs",
        "//app/service/live/rtc/common:all-srcs",
        "//app/service/live/rtc/internal/conf:all-srcs",
        "//app/service/live/rtc/internal/dao:all-srcs",
        "//app/service/live/rtc/internal/model:all-srcs",
        "//app/service/live/rtc/internal/server/grpc:all-srcs",
        "//app/service/live/rtc/internal/server/http:all-srcs",
        "//app/service/live/rtc/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
