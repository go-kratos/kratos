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
        "//app/job/main/tv/service/pgc:all-srcs",
        "//app/job/main/tv/service/report:all-srcs",
        "//app/job/main/tv/service/ugc:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
