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
        "//app/job/main/member-cache/cmd:all-srcs",
        "//app/job/main/member-cache/conf:all-srcs",
        "//app/job/main/member-cache/dao:all-srcs",
        "//app/job/main/member-cache/model:all-srcs",
        "//app/job/main/member-cache/server/http:all-srcs",
        "//app/job/main/member-cache/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
