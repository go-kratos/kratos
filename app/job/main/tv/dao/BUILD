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
        "//app/job/main/tv/dao/app:all-srcs",
        "//app/job/main/tv/dao/archive:all-srcs",
        "//app/job/main/tv/dao/cms:all-srcs",
        "//app/job/main/tv/dao/ftp:all-srcs",
        "//app/job/main/tv/dao/lic:all-srcs",
        "//app/job/main/tv/dao/playurl:all-srcs",
        "//app/job/main/tv/dao/report:all-srcs",
        "//app/job/main/tv/dao/ugc:all-srcs",
        "//app/job/main/tv/dao/upper:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
