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
        "//app/service/live/live_riskcontrol/api/grpc/v1:all-srcs",
        "//app/service/live/live_riskcontrol/cmd:all-srcs",
        "//app/service/live/live_riskcontrol/internal/conf:all-srcs",
        "//app/service/live/live_riskcontrol/internal/dao:all-srcs",
        "//app/service/live/live_riskcontrol/internal/model:all-srcs",
        "//app/service/live/live_riskcontrol/internal/server/http:all-srcs",
        "//app/service/live/live_riskcontrol/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
