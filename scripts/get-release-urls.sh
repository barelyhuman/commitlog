
#!/usr/bin/env bash

curl -sL https://api.github.com/repos/barelyhuman/commitlog/releases | jq '[ .[0].assets[] | {name:.name,url:.url} ]' > api_urls.json

cat >./docs/download.md<<EOF
<!-- meta -->
<title>
    commitlog | downloads
</title>
<meta name="description" content="commits to changelog generator">
<!-- meta end -->

### [commitlog](/)

[Home](/) [Manual](/manual) [Download](/download) [API](/api)

#### Downloads

EOF

jq -c '.[]' api_urls.json | while read i; do
    url=$(echo $i | jq '.url' --raw-output)
    name=$(echo $i | jq '.name' --raw-output)
    bdurl=$(curl -sL $url | jq '.browser_download_url')
    echo "[$name]($bdurl)    " >> ./docs/download.md
done