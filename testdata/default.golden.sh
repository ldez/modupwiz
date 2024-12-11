#!/bin/sh -e

git fetch --multiple upstream origin
git reset --hard upstream/master

# https://github.com/stretchr/testify/compare/v1.9.0...v1.10.0
echo "Update github.com/stretchr/testify@v1.10.0 "
go get github.com/stretchr/testify@v1.10.0
go mod tidy
git add .; git commit -m "github.com/stretchr/testify@v1.10.0"

# https://github.com/golang/text/compare/v0.19.0...v0.21.0
echo "Update golang.org/x/text@v0.21.0 "
go get golang.org/x/text@v0.21.0
go mod tidy
git add .; git commit -m "golang.org/x/text@v0.21.0"

# https://github.com/golang/crypto/compare/v0.28.0...v0.31.0
# https://github.com/golang/net/compare/v0.30.0...v0.32.0
echo "Update golang.org/x/crypto@v0.31.0 golang.org/x/net@v0.32.0 "
go get golang.org/x/crypto@v0.31.0 golang.org/x/net@v0.32.0
go mod tidy
git add .; git commit -m "Dependency group" -m "golang.org/x/crypto@v0.31.0" -m "golang.org/x/net@v0.32.0"

# https://github.com/golang/time/compare/v0.7.0...v0.8.0
echo "Update golang.org/x/time@v0.8.0 "
go get golang.org/x/time@v0.8.0
go mod tidy
git add .; git commit -m "golang.org/x/time@v0.8.0"

# https://github.com/cloudflare/cloudflare-go/compare/v0.108.0...v0.111.0
echo "Update github.com/cloudflare/cloudflare-go@v0.111.0 "
go get github.com/cloudflare/cloudflare-go@v0.111.0
go mod tidy
git add .; git commit -m "github.com/cloudflare/cloudflare-go@v0.111.0"

# https://github.com/nrdcg/freemyip/compare/v0.2.0...v0.3.0
echo "Update github.com/nrdcg/freemyip@v0.3.0 "
go get github.com/nrdcg/freemyip@v0.3.0
go mod tidy
git add .; git commit -m "github.com/nrdcg/freemyip@v0.3.0"

# https://github.com/yandex-cloud/go-genproto/compare/76a0cfc1a773...07e4a676108b
echo "Update github.com/yandex-cloud/go-genproto@v0.0.0-20241206133605-07e4a676108b "
go get github.com/yandex-cloud/go-genproto@v0.0.0-20241206133605-07e4a676108b
go mod tidy
git add .; git commit -m "github.com/yandex-cloud/go-genproto@v0.0.0-20241206133605-07e4a676108b"

# https://github.com/yandex-cloud/go-sdk/compare/947cf519f6bd...6c3760d17eea
echo "Update github.com/yandex-cloud/go-sdk@v0.0.0-20241206142255-6c3760d17eea "
go get github.com/yandex-cloud/go-sdk@v0.0.0-20241206142255-6c3760d17eea
go mod tidy
git add .; git commit -m "github.com/yandex-cloud/go-sdk@v0.0.0-20241206142255-6c3760d17eea"

# https://github.com/aliyun/alibaba-cloud-sdk-go/compare/v1.63.47...v1.63.66
echo "Update github.com/aliyun/alibaba-cloud-sdk-go@v1.63.66 "
go get github.com/aliyun/alibaba-cloud-sdk-go@v1.63.66
go mod tidy
git add .; git commit -m "github.com/aliyun/alibaba-cloud-sdk-go@v1.63.66"

# https://github.com/aws/aws-sdk-go-v2/compare/v1.32.3...v1.32.6
echo "Update github.com/aws/aws-sdk-go-v2@v1.32.6 "
go get github.com/aws/aws-sdk-go-v2@v1.32.6
go mod tidy
git add .; git commit -m "github.com/aws/aws-sdk-go-v2@v1.32.6"

# https://github.com/aws/aws-sdk-go-v2/compare/service/sts/v1.32.3...service/sts/v1.33.2
echo "Update github.com/aws/aws-sdk-go-v2/service/sts@v1.33.2 "
go get github.com/aws/aws-sdk-go-v2/service/sts@v1.33.2
go mod tidy
git add .; git commit -m "github.com/aws/aws-sdk-go-v2/service/sts@v1.33.2"

# https://github.com/aws/aws-sdk-go-v2/compare/credentials/v1.17.42...credentials/v1.17.47
echo "Update github.com/aws/aws-sdk-go-v2/credentials@v1.17.47 "
go get github.com/aws/aws-sdk-go-v2/credentials@v1.17.47
go mod tidy
git add .; git commit -m "github.com/aws/aws-sdk-go-v2/credentials@v1.17.47"

# https://github.com/aws/aws-sdk-go-v2/compare/config/v1.28.1...config/v1.28.6
echo "Update github.com/aws/aws-sdk-go-v2/config@v1.28.6 "
go get github.com/aws/aws-sdk-go-v2/config@v1.28.6
go mod tidy
git add .; git commit -m "github.com/aws/aws-sdk-go-v2/config@v1.28.6"

# https://github.com/aws/aws-sdk-go-v2/compare/service/s3/v1.66.2...service/s3/v1.71.0
echo "Update github.com/aws/aws-sdk-go-v2/service/s3@v1.71.0 "
go get github.com/aws/aws-sdk-go-v2/service/s3@v1.71.0
go mod tidy
git add .; git commit -m "github.com/aws/aws-sdk-go-v2/service/s3@v1.71.0"

# https://github.com/golang/oauth2/compare/v0.23.0...v0.24.0
echo "Update golang.org/x/oauth2@v0.24.0 "
go get golang.org/x/oauth2@v0.24.0
go mod tidy
git add .; git commit -m "golang.org/x/oauth2@v0.24.0"

# https://github.com/aws/aws-sdk-go-v2/compare/service/lightsail/v1.42.3...service/lightsail/v1.42.7
echo "Update github.com/aws/aws-sdk-go-v2/service/lightsail@v1.42.7 "
go get github.com/aws/aws-sdk-go-v2/service/lightsail@v1.42.7
go mod tidy
git add .; git commit -m "github.com/aws/aws-sdk-go-v2/service/lightsail@v1.42.7"

# https://github.com/aws/aws-sdk-go-v2/compare/service/route53/v1.46.0...service/route53/v1.46.3
echo "Update github.com/aws/aws-sdk-go-v2/service/route53@v1.46.3 "
go get github.com/aws/aws-sdk-go-v2/service/route53@v1.46.3
go mod tidy
git add .; git commit -m "github.com/aws/aws-sdk-go-v2/service/route53@v1.46.3"

# https://github.com/civo/civogo/compare/v0.3.11...v0.3.89
echo "Update github.com/civo/civogo@v0.3.89 "
go get github.com/civo/civogo@v0.3.89
go mod tidy
git add .; git commit -m "github.com/civo/civogo@v0.3.89"

# https://github.com/huaweicloud/huaweicloud-sdk-go-v3/compare/v0.1.120...v0.1.126
echo "Update github.com/huaweicloud/huaweicloud-sdk-go-v3@v0.1.126 "
go get github.com/huaweicloud/huaweicloud-sdk-go-v3@v0.1.126
go mod tidy
git add .; git commit -m "github.com/huaweicloud/huaweicloud-sdk-go-v3@v0.1.126"

# https://github.com/linode/linodego/compare/v1.42.0...v1.43.0
echo "Update github.com/linode/linodego@v1.43.0 "
go get github.com/linode/linodego@v1.43.0
go mod tidy
git add .; git commit -m "github.com/linode/linodego@v1.43.0"

# https://github.com/nrdcg/desec/compare/v0.8.0...v0.10.0
echo "Update github.com/nrdcg/desec@v0.10.0 "
go get github.com/nrdcg/desec@v0.10.0
go mod tidy
git add .; git commit -m "github.com/nrdcg/desec@v0.10.0"

# https://github.com/oracle/oci-go-sdk/compare/v65.77.1...v65.80.0
echo "Update github.com/oracle/oci-go-sdk/v65@v65.80.0 "
go get github.com/oracle/oci-go-sdk/v65@v65.80.0
go mod tidy
git add .; git commit -m "github.com/oracle/oci-go-sdk/v65@v65.80.0"

# https://github.com/selectel/go-selvpcclient/compare/v3.1.1...v3.2.1
echo "Update github.com/selectel/go-selvpcclient/v3@v3.2.1 "
go get github.com/selectel/go-selvpcclient/v3@v3.2.1
go mod tidy
git add .; git commit -m "github.com/selectel/go-selvpcclient/v3@v3.2.1"

# https://github.com/tencentcloud/tencentcloud-sdk-go/compare/tencentcloud/common/v1.0.1034...tencentcloud/common/v1.0.1058
echo "Update github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common@v1.0.1058 "
go get github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common@v1.0.1058
go mod tidy
git add .; git commit -m "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common@v1.0.1058"

# https://github.com/tencentcloud/tencentcloud-sdk-go/compare/tencentcloud/dnspod/v1.0.1034...tencentcloud/dnspod/v1.0.1058
echo "Update github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod@v1.0.1058 "
go get github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod@v1.0.1058
go mod tidy
git add .; git commit -m "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod@v1.0.1058"

# https://github.com/volcengine/volc-sdk-golang/compare/v1.0.183...v1.0.188
echo "Update github.com/volcengine/volc-sdk-golang@v1.0.188 "
go get github.com/volcengine/volc-sdk-golang@v1.0.188
go mod tidy
git add .; git commit -m "github.com/volcengine/volc-sdk-golang@v1.0.188"

# https://github.com/vultr/govultr/compare/v3.9.1...v3.12.0
echo "Update github.com/vultr/govultr/v3@v3.12.0 "
go get github.com/vultr/govultr/v3@v3.12.0
go mod tidy
git add .; git commit -m "github.com/vultr/govultr/v3@v3.12.0"

# https://github.com/googleapis/google-api-go-client/compare/v0.204.0...v0.211.0
echo "Update google.golang.org/api@v0.211.0 "
go get google.golang.org/api@v0.211.0
go mod tidy
git add .; git commit -m "google.golang.org/api@v0.211.0"

