#!/usr/bin/env bash
#替换当前目录下文件中Id,Url这种不合规的字符

replacements="Id Url Sku Uid"

for i in $replacements; do
    r=`echo $i|tr 'a-z' 'A-Z'`
    find . -name "*.go"|xargs sed -ibak -E 's/'$i'(s?)([[:>:]]|[A-Z_])/'$r'\1\2/g'
done

find . -name "*.gobak"|xargs rm