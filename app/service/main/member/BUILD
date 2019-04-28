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
        "//app/service/main/member/api:all-srcs",
        "//app/service/main/member/cmd:all-srcs",
        "//app/service/main/member/conf:all-srcs",
        "//app/service/main/member/dao:all-srcs",
        "//app/service/main/member/model:all-srcs",
        "//app/service/main/member/server:all-srcs",
        "//app/service/main/member/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
