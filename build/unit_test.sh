#!/bin/bash

CI_SERVER_URL="http://git.bilibili.co"
CI_UATSVEN_URL="http://uat-sven.bilibili.co"
CI_PRIVATE_TOKEN="WVYk-ezyXKq-C82v-1Bi"
CI_PROJECT_ID="682"
CI_COMMIT_SHA=${PULL_PULL_SHA}

exitCode=0

# get packages
dirs=(dao)
declare packages
function GetPackages(){
    reg="library"
    length=${#dirs[@]}
    for((i=0;i<length;i++))
    do
        reg+="|$1/"${dirs[i]}"(/.*)*"

    done
    for value in `find $1 -type d |grep -E ${reg}`
    do
        len=`ls ${value}| grep .go | wc -l`
        if [[ ${len} -gt 0 ]];then
            packages+="go-common/"${value}" "
        fi
    done
}
# upload data to apm
# $1: SvenURL
# $2: file result.out path
function Upload () {
    if [[ ! -f "$2/result.out" ]] || [[ ! -f "$2/cover.html" ]] || [[ ! -f "$2/coverage.dat" ]]; then
        echo "==================================WARNING!======================================"
        echo "No test found!~ 请完善如下路径测试用例： ${pkg} "
        exit 1
    fi
    json=$(curl $1 -H "Content-type: multipart/form-data" -F "html_file=@$2/cover.html" -F "report_file=@$2/result.out" -F "data_file=@$2/coverage.dat")
    if [[ "${json}" = "" ]]; then
        echo "shell.Upload curl $1 fail"
        exit 1
    fi
    msg=$(echo ${json} | jq -r '.message')
    data=$(echo ${json} | jq -r '.data')
    code=$(echo ${json} | jq -r '.code')
    if [[ "${data}" = "" ]]; then
        echo "shell.Upload curl $1 fail,data return null"
        exit 1
    fi
    echo "=============================================================================="
    if [[ ${code} -ne 0 ]]; then
        echo -e "返回 message(${msg})"
        echo -e "返回 data(${data})\n\n"
    fi
    return ${code}
}

# GoTest execute go test and go tool
# $1: pkg
function GoTest(){
    go test -v $1 -coverprofile=cover.out -covermode=set -convey-json -timeout=60s > result.out
    go tool cover -html=cover.out -o cover.html
}

# BazelTest execute bazel coverage and go tool
# $1: pkg
function BazelTest(){
    pkg=${1//go-common//}":go_default_test"
    path=${1//go-common\//}

    bazel coverage --instrumentation_filter="//${path}[:]" --test_env=DEPLOY_ENV=uat --test_timeout=60 --test_env=APP_ID=bazel.test --test_output=all --cache_test_results=no --test_arg=-convey-json ${pkg} > result.out
    if [[ ! -s result.out ]]; then
        echo "==================================WARNING!======================================"
        echo "No test case found,请完善如下路径测试用例： ${pkg} "
        exit 1
    else
        echo $?
        cat bazel-out/k8-fastbuild/testlogs/${path}/go_default_test/coverage.dat | grep -v "/monkey.go" > coverage.dat
        go tool cover -html=coverage.dat -o cover.html
    fi
}

# UTLint check the *_test.go files in the pkg
# $1: pkg
function UTLint()
{
    path=${1//go-common\//}
    declare -i numCase=0
    declare -i numAssertion=0
    files=$(ls ${path} | grep -E "(.*)_test\.go")
    if [[ ${#files} -eq 0 ]];then
        echo "shell.UTLint no *_test.go files in pkg:$1"
        exit 1
    fi
    for file in ${files}
    do
        numCase+=`grep -c -E "^func Test(.+)\(t \*testing\.T\) \{$" ${path}/${file}`
        numAssertion+=`grep -c -E "^(.*)So\((.+)\)$" ${path}/${file}`
    done
    if [[ ${numCase} -eq 0 || ${numAssertion} -eq 0 ]];then
        echo -e "shell.UTLint no test case or assertion in pkg:$1"
        exit 1
    fi
    echo "shell.UTLint pkg:$1 succeeded"
}

# upload path to apm
# $1: SvenURL
# $2: file result.out path
function UpPath() {
    curl $1 -H "Content-type: multipart/form-data" -F "path_file=@$2/path.out"
}
function ReadDir(){
    # get go-common/app all dir path
    PathDirs=`find app -maxdepth 3 -type d`
    value=""
    for dir in ${PathDirs}
    do
        if [[ -d "$dir" ]];then
            for file in `find ${dir} -maxdepth 1 -type f |grep "CONTRIBUTORS.md"`
            do
                owner=""
                substr=${dir#*"go-common"}
                while read line
                do
                    if [[ "${line}" = "# Owner" ]];then
                        continue
                    elif [[ "${line}" = "" ]]|| [[ "${line}" = "#"* ]];then
                        break
                    else
                        owner+="${line},"
                    fi
                done < ${file}
                value+="{\"path\":\"go-common${substr}\",\"owner\":\"${owner%,}\"},"
            done
        fi
    done
    # delete "," at the end of value
    value=${value%,}
    echo "[${value}]" > path.out
}

# start work
function Start(){
    GetPackages $1
    if [[ ${packages} = "" ]]; then
        echo "shell.Start no change packages"
        exit 0
    fi
    #Get gitlab result
    gitMergeRequestUrl="${CI_SERVER_URL}/api/v4/projects/${CI_PROJECT_ID}/repository/commits/${CI_COMMIT_SHA}/merge_requests?private_token=${CI_PRIVATE_TOKEN}"
    gitCommitUrl="${CI_SERVER_URL}/api/v4/projects/${CI_PROJECT_ID}/repository/commits/${CI_COMMIT_SHA}/statuses?private_token=${CI_PRIVATE_TOKEN}"
    mergeJson=$(curl -s ${gitMergeRequestUrl})
    commitJson=$(curl -s ${gitCommitUrl})
    if [[ "${mergeJson}" != "[]" ]] && [[ "${commitJson}" != "[]" ]]; then
        merge_id=$(echo ${mergeJson} | jq -r '.[0].iid')
        exitCode=$?
        if [[ ${exitCode}  -ne 0 ]]; then
            echo "shell.Start curl ${gitMergeRequestUrl%=*}=*** error .return(${mergeJson})"
            exit 1
        fi
        username=$(echo ${mergeJson} | jq -r '.[0].author.username')
        authorname=$(echo ${commitJson} | jq -r '.[0].author.username')
    else
        echo "Test not run, maybe you should try create a merge request first!"
        exit 0
    fi
    #Magic time
    Magic
    #Normal process
    for pkg in ${packages}
    do
        svenUrl="${CI_UATSVEN_URL}/x/admin/apm/ut/upload?merge_id=${merge_id}&username=${username}&author=${authorname}&commit_id=${CI_COMMIT_SHA}&pkg=${pkg}"
        echo "shell.Start ut lint pkg:${pkg}"
        UTLint "${pkg}"
        echo "shell.Start Go bazel test pkg:${pkg}"
        BazelTest "${pkg}"
        Upload ${svenUrl} $(pwd)
        exitCode=$?
        if [[ ${exitCode} -ne 0 ]]; then
            echo "shell.Start upload fail, status(${exitCode})"
            exit 1
        fi
    done
    # upload all dirs
    ReadDir
    pathUrl="${CI_UATSVEN_URL}/x/admin/apm/ut/upload/app"
    UpPath ${pathUrl} $(pwd)
    echo "UpPath has finshed......  $(pwd)"
    return 0
}

# Check determine whether the standard is up to standard
#$1: commit_id
function Check(){
    curl "${CI_UATSVEN_URL}/x/admin/apm/ut/git/report?project_id=${CI_PROJECT_ID}&merge_id=${merge_id}&commit_id=$1"
    checkURL="${CI_UATSVEN_URL}/x/admin/apm/ut/check?commit_id=$1"
    json=$(curl -s ${checkURL})
    code=$(echo ${json} | jq -r '.code')
    if [[ ${code} -ne 0 ]]; then
        echo -e "curl ${checkURL} response(${json})"
        exit 1
    fi
    package=$(echo ${json} | jq -r '.data.package')
    coverage=$(echo ${json} | jq -r '.data.coverage')
    passRate=$(echo ${json} | jq -r '.data.pass_rate')
    standard=$(echo ${json} | jq -r '.data.standard')
    increase=$(echo ${json} | jq -r '.data.increase')
    tyrant=$(echo ${json} | jq -r '.data.tyrant')
    lastCID=$(echo ${json} | jq -r '.data.last_cid')
    if ${tyrant}; then
        echo -e "\t续命失败!\n\t大佬，本次执行结果未达标哦(灬ꈍ ꈍ灬)，请再次优化ut重新提交🆙"
        echo -e "\t---------------------------------------------------------------------"
        printf "\t%-14s %-14s %-14s %-14s\n" "本次覆盖率(%)" "本次通过率(%)" "本次增长量(%)" 执行pkg
        printf "\t%-13.2f %-13.2f %-13.2f %-12s\n" ${coverage} ${passRate} ${increase} ${package}
        echo -e "\t(达标标准：覆盖率>=${standard} && 通过率=100% && 同比当前package历史最高覆盖率的增长率>=0)"
        echo -e "\t---------------------------------------------------------------------"
        exitCode=1
    else
        echo -e "\t恭喜你，续命成功，可以请求MR了"
    fi
}

# Magic ignore method Check()
function Magic(){
    url="http://git.bilibili.co/api/v4/projects/${CI_PROJECT_ID}/merge_requests/${merge_id}/notes?private_token=${CI_PRIVATE_TOKEN}"
    json=$(curl -s ${url})
    for comment in $(echo ${json} | jq -r '.[].body')
    do
        if [[ ${comment} == "+skiput" ]]; then
            exit 0
        fi
    done
}


# run
Start $1
echo -e "【我们万众一心】："
Check ${CI_COMMIT_SHA}
echo -e "本次执行详细结果查询地址请访问：http://sven.bilibili.co/#/ut?merge_id=${merge_id}&&pn=1&ps=20"
if [[ ${exitCode} -ne 0 ]]; then
    echo -e "执行失败！！！请解决问题后再次提交。具体请参考：http://info.bilibili.co/pages/viewpage.action?pageId=9841745"
    exit 1
else
    echo -e "执行成功."
    exit 0
fi
