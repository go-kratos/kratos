package(default_visibility = ["//visibility:public"])

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
        "//app/interface/live/push-live/cmd:all-srcs",
        "//app/interface/live/push-live/conf:all-srcs",
        "//app/interface/live/push-live/dao:all-srcs",
        "//app/interface/live/push-live/http:all-srcs",
        "//app/interface/live/push-live/model:all-srcs",
        "//app/interface/live/push-live/service:all-srcs",
    ],
    tags = ["automanaged"],
)
