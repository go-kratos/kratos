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
        "//app/service/main/antispam/cmd:all-srcs",
        "//app/service/main/antispam/conf:all-srcs",
        "//app/service/main/antispam/dao:all-srcs",
        "//app/service/main/antispam/extern:all-srcs",
        "//app/service/main/antispam/http:all-srcs",
        "//app/service/main/antispam/model:all-srcs",
        "//app/service/main/antispam/rpc/client:all-srcs",
        "//app/service/main/antispam/rpc/server:all-srcs",
        "//app/service/main/antispam/service:all-srcs",
        "//app/service/main/antispam/util:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
