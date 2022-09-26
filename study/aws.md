
# AWS

## IAM(Identity and Access Management)
웹서비스를 안전하게 AWS 리소스를 컨트롤 할 수 있게 해주는 역할
AWS 리소스를 사용하도록 인증되고 승인된 사람을 제어하는데 사용할 수 있다.

User groups를 사용하면 사용자 그룹별로 aws 리소스 권한을 부여할 수 있음(admin, developer, tester 나눠서 사용)

## AWS DB 설정
aws RDS에 db가 올라가고 Makefile에 작성한 migrate 스크립트를 localhost에서 RDS의 url과 password로 변경하게 된다.
실제 운영환경에서는 app.env에서 DB_SOURCE를 RDS의 DB url과 password로 변경해야한다.
TOKEN 대칭 키도 마찬가지인데
환경변수를 test용이랑 real production용이랑 적용을 다르게 해야한다.

docker를 빌드하고 ECR에 컨테이너가 올라가기 전에 환경변수를 셋팅해주면 되는데 이 환경변수들은 보안적으로 중요한 정보를 담고 있으므로 github-repo에 올리면 안된다.
좋은 방법은 AWS secret manager service를 사용하는 것이다.

AWS secret manager service에 저장한 값을 app.env에서 검색하게 만들려면 aws cli를 설치하여 AWS에 접근할 수 있게 만들자.
aws-cli를 설치했으면 configure를 설정해줘야한다.

```bash
aws configure
```
access_key랑 secret key를 등록해야하는데 AWS IAM에서 user-security_credentials에서 access_key를 만들어서 등록해주면 된다.


```bash
aws secretsmanager get-secret-value --secret-id {FRIENDLY_NAME OR ARN} --query SecretString --output text
```
위 스크립트를 사용하여 aws secret manager에 등록된 환경변수를 json으로 받아올 수 있다. 스크립트 실행 결과는 raw text로 나오기 때문에 json 포멧에 맞춰 변경해줘야한다.
이를 위해 **jq**를 설치해줘야한다.
> brew install jq

환경변수 뽑는 스크립트를 파이프로 jq에 전달하면 json 포멧에 맞춰 잘 나온다.
```bash

```bash
aws secretsmanager get-secret-value --secret-id {FRIENDLY_NAME OR ARN} --query SecretString --output text | jq 'to_entries|map("\(.key)=\(.value)")|.[]' -r > app.env
```
위 스크립트로 aws-cli로 뽑은 aws secret 환경변수 값들을 동적으로 app.env에 적용해주면 된다.

## kubernetes component

+ master node
  + worker node들을 관리하기 위함
  + 프론트 엔드에서 요청을 처리하기 위한 API 서버로 master node가 될 수 있다.
  + etcd는 모든 클러스터 데이터를 key:value로 저장한다.
  + scheduler는 아직 노드에 할당되지 않은 새 pod를 감시하다가 실행하기 위해 노드를 선택한다.
  + control manager는 다음 컨트롤러들을 모아둔거다.
    + node controller
    + job controller
    + end-point controller
    + service account & token controller
  + cloud controller manager는 cloud provider api와 통신하기 위한 매니저다.
    + node-controller
    + route controller
    + service controller
    + 
+ worker node
  + 각 worker node에는 kubelet agent가 존재
  + kubelet은 pod를 관리함 
  + kube-proxy
    + 쿠버네티스 네트워크간 통신

## EKS
쿠버네티스 컴포넌트가 복합적이므로 쉽게 관리하기 위한게 AWS EKS다.
마스터 노드는 EKS에서 자동으로 생성된다.
설정 해야할 일은
+ 워커노드를 EKS에 추가해주고 관련 셋팅들만 해주면된다.

aws-cli를 통해 로컬 터미널에서 클라우드로 접속할 수 있다.
IAM에서 EKS에 필요한 ROLE을 설정해주고 EKS에 적용해준다.

aws-cli의 credentials에서 이미 다른 AWS 서비스(github에서 ECR 연동을 위해 authorized된 access_key, secret_access_key)를 사용하고 있으면 ***~/.aws/credentials*** 파일에 EKS용으로 새로 key를 만들어서 추가해주면 된다.
```bash
[default]
aws_access_key_id = [ACCESS_KEY]
aws_secret_access_key_id = [SECRET_ACCESS_KEY]

[github]
aws_access_key_id = [ACCESS_KEY]
aws_secret_access_key_id = [SECRET_ACCESS_KEY]
```
위 예시처럼 aws-cli를 github 관련 용으로 쓸 때는 ***AWS_PROFILE*** 환경변수를 조정해주면 된다.
> export AWS_PROFILE=github

그렇다면 _github_ credential로  클러스터에 접근하려면 어떻게 해야 할까?
https://aws.amazon.com/ko/premiumsupport/knowledge-center/amazon-eks-cluster-access/
이 내용에 따르면 aws-auth.yaml에 접근하려는 사용자의 arn을 등록해주면 된다.
```yaml
#aws-auth.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: aws-auth
  namespace: kube-system
data:
  mapUsers: |
    - userarn: {IAM 사용자의 ARN}
      username: {IAM에 등록된 사용자 이름}
      groups:
        - system:masters
```

> kubectl apply -f aws-auth.yaml

제대로 작동하는지 확인하기 위해서 AWS_PROFILE을 default, github 번갈아가며 ***kubectl cluster-info*** 로 테스트 했을 때 동일한 결과가 나와야한다.

## k8s 모니터링
> brew install k9s

> k9s


## deployment와 service
만들어진 API 웹서버는 k8s의 deployment로 배포해야하고 (이 프로젝트 루트의 /eks/deployment.yaml)
외부에서 external-ip로 접근하기 위해서 k8s의 service를 만들어 load-balancer로 배포해야한다.
yaml 작성 방법은 https://kubernetes.io/docs/concepts/services-networking/service/ 공식문서를 참고하자

## ingress와 ingress-controller
ingress는 외부에서 k8s 클러스터로 접근을 관리하는 API 오브젝트이며 HTTP/HTTPS를 관리한다.
+ 부하분산, SSL 종료, 명칭 기반 가상 호스팅을 제공한다.
+ 트래픽 라우팅은 ingress 리소스에 정의된 규칙에 의해 컨트롤된다.
+ nginx로 라우팅을 관리함

그럼 이게 왜 필요한가?
요청을 서비스 별로 나눠서 부하를 분산시킬 수 있기 때문이다.
가령 회원가입과 로그인 요청을 k8s의 한 서비스로 받아 냈다면 
ingress를 사용한다면 서비스를 회원가입 서비스와 로그인 서비스로 나눠서 받을 수 있다.

## Cert-manager
인증서를 생성하고 시간이 지나면 만료된다.
cert-manager는 자동으로 인증서를 새로 발급해주는 역할을 한다.
ingress.yaml에 annotation으로 맞물려주면 된다.