DB migration

# Schema Migration
관계형 데이터베이스 스키마에 대한 버전 제어, 증분 및 되돌릴 수 있는 변경의 관리를 나타냅니다. 데이터베이스의 스키마를 최신 또는 이전 버전으로 업데이트하거나 되돌릴 필요가 있을 때마다 데이터베이스에서 스키마 마이그레이션이 수행됩니다.

migration up : 스키마 버전 올리는 것
migration down : 스키마 버전 내리는 것
# Data Migration
DB 버전 업그레이드, 노후화된 노드의 변경 및 DB 교체 등의 이유로 운영 중인 DB를 변경해야 하는 경우가 종종 발생한다. 이 때 기존 DB에 저장되어 있던 데이터를 신규 DB로 옮기는 과정이 필요한데 이를 데이터 마이그레이션이라고 한다.

# golang-migrate
db migration을 위한 유틸리티
> brew install golang-migrate

migrate up/down할 sql 스켈레톤 파일을 생성해줌
> migrate create -ext sql -dir {target_dir} -seq init_schema

## migration up
요구사항에 의해 DB 스키마 변경된 경우 기존에 생성한 스켈레톤 파일에 덮어씌우는게 아니라 새로운 파일을 만들어서 관리한다.
> migrate create -ext sql -dir {target_dir} -seq something_schema_changed

# postgres

## db 생성
```bash
createdb --username=root --owner=root {db_name}
```

## db 접속
```bash
psql {db_name}
```

## db 삭제
```bash
dropdb {db_name}
```

# SQLC
sqlc는 type-safe한 코드를 생성해준다.
매우 빠르며 사용하기 쉽다.(gorm은 느리며 복잡한 쿼리 작성시 코드 작성이 쉽지 않다. 예제도 별로 없고..)
postgresql을 잘 지원해주기 때문에 

5.44부터


# Transaction Isolation Level
동시 트랜잭션은 서로 영향을 미치지 않아야하는 데이터 베이스 성질을 Isolation이라 한다.

postgres에선 다음 명령어로 isolation level을 확인할 수 있다.
> show transaction isolation level

isolation level은 4단계가 있다.

https://www.postgresql.org/docs/current/transaction-iso.html

## 명심할 것
- isolation 레벨이 높을 수록 error, timeout, deadlock이 발생할 수 있다.
- DBMS마다 isolation 구현이 다르니 공식 문서를 꼭 참조하자.

# PK와 UNIQUE INDEX

## 아래는 친절한 SQL 튜닝 저자인 조시형씨가 2004년에 적은 글
안녕하십니까? 엔코아 정보컨설팅에 근무하는 컨설턴트 조시형입니다.

우선 Primary Key와 Unique Index의 차이점을 설명하는 것은 부적절하다는 말씀을 드리고 싶습니다. 둘간의 상관관계를 설명하는 것이 맞는 개념입니다. 많은 개발자들이 PK는 왠지 부하를 준다는 잘못된 선입견을 가지고 있고 따라서 PK 대신 Unique Index를 사용하는 것으로 알고 있는데 매우 그릇된 관행(?)이라고 생각합니다. 따라서 앞에서 다른 분들이 좋은 설명 많이 해 주셨지만 부연해서 설명을 드리도록 하겠습니다.

Primary Key라고 하는 것은 논리적인 개념입니다. Primary Key는 해당 컬럼이 그 테이블의 식별자임을 나타내는 것으로서, 자신과 다른 레코드가 서로 다른 인스턴스임을 확인할 수 있게 해 주는 역할을 합니다. 즉, 해당 그 레코드의 존재자체인 것이지요.
원래 사람이름이라는 것이 '나'와 다른 사람을 식별하기 위해 사용하는 것인데, 다른 사람과 중복될 수 있으므로 나를 식별할 수 있는 속성으로서 주민등록번호라는 것을 대신 사용합니다. 따라서 주민등록번호는 나와 별개가 아닌 나의 존재 그 자체입니다. 철학적으로는 맞지 않는 설명이겠지만 적어도 ‘데이터의 세계’에서는 그렇습니다.

반면 PK constraint는 물리적인 개념입니다. "이 컬럼(들)은 다른 레코드와 구분짓는 식별자 역할을 하는 중요한 컬럼이므로 데이터는 중복을 허용해서는 안 되고(unique), null값을 허용해서도 안돼(not null)"라고 DB에게 정보를 주는 것입니다. primary key constraint를 설정하면 unique index와 not null constraint가 자동적으로 생성되는 이유도 여기에 있지요. 
Primary Key가 논리적인 개념이며, PK Constraint 와 not null 등의 제약조건과 unique index 들이 물리적인 개념이라고 하는게 맞겠군요.

특히 인덱스는 PK 컬럼의 Unique성을 보장하기 위해 매우 필수적인 도구인데, 만약 인덱스 없이 해당 컬럼값이 중복되지 않도록 할 수 있는 방법을 생각해 보시기 바랍니다. 잘 떠오르지 않을 것입니다. 결론적으로 말씀드리면, 인덱스와 PK의 상관관계에 있어서 Unique이든 Non-Unique이든 Index라는 놈은 PK 컬럼의 Unique성을 보장하기 위해 사용하는 하나의 도구(Tool)에 지나지 않습니다.
어떤 의미에서만 본다면, 어느 한 컬럼에 Not Null Constraint를 주고, 그 컬럼에 Unique Index를 생성하였다면 primary key와 다를 것이 하나도 없어 보입니다.

하지만 primary key는 데이터베이스와 사용자 입장에서 매우 중요한 정보 역할을 하기 때문에 중요하다고 말씀드리고 싶습니다.
앞서 말씀드렸듯이, 원래 primary key는 “이 컬럼(들)이 테이블의 식별자(identifier)이므로 중복을 허용해서는 안 되고, null값을 허용해서도 안 된다”라는 의미론적(semantic) 인 의미에서 정의하는 것이며, DBMS는 이를 효과적으로 처리하기 위해서 index를 자동생성해서 사용하고 not null constraint를 정의하는 것입니다.

참고로 primary key를 위해서 반드시 unique index가 필요한 것은 아닙니다. non-unique index만 있더라도 새로운 값이 들어올 때 중복값이 있는지 체크하는 데에는 전혀 문제가 없으므로 기존에 이미 non-unique 인덱스가 정의되어 있는 상황이었다면 그 인덱스를 그대로 사용합니다. DW 시스템에서 대량의 데이터 로딩시 속도를 빠르게 하기 위해 PK Constraint를 일시적으로 Disable 시키는 경우가 있는데 이렇게 되면 Unique Index도 동시에 제거되므로 인덱스를 다시 생성해야 하는 부담이 생깁니다. 이러한 부담을 덜기 위해 의도적으로 non-unique index를 생성하는 경우가 있는데 이에 대해서는 더 깊이 언급하지 않겠습니다.

primary constraint key 를 생성하면 자동적으로 unique index가 생기는데, 주의할점은 primary constraint key를 해제할 때 unique index 도 자동적으로 삭제된다, 따라서 연기가능한 제약조건(맞나?) 등을 이용할 수 있다. 라는 부분이었습니다.
그때는 사람이 실수나, 다른 이유로 제약조건을 삭제할 때 인덱스까지 없어지므로 갑자기 느려지는 상황이 올 수 있다. 라고 말을 했었는데요, 그런 경우와 함께 DW 시스템 (OLAP) 에서 데이터 로딩 속도를 빠르게 하기 위해 PK Constraint 를 일시적으로 Disable 하는 경우가 있어서, 이것땜에 인덱스를 다시 생성하는 경우가 있다. 그래서 의도적으로 non-unique index 를 생성하는 경우가 있다.. 라고 하네요.

하여튼 primary key, foreign key, not null 등과 같은 integrity constraint 정보들은 plan을 생성하고, query rewrite(주로 DW에서 많이 사용되는 기능임)를 수행할 때, 그리고 기타 여러가지 용도로 데이터베이스에 의해 사용되어집니다.
그리고 OLAP Tool과 같이 데이터베이스에 접근하는 여러 tool들이 동적으로 ad-hoc query를 생성할 때 이 정보들을 활용합니다.

optimizer와 tool 입장에서 뿐만 아니라 데이터베이스를 사용하는 사용자 입장에서도 실제 document 역할을 하게 되므로 의미가 있습니다. ER Diagram을 보지 않고 data dictionary에 있는 정보만을 보고도 그 테이블의 식별자(Identifier)가 어떤 컬럼으로 구성되어 있는지 쉽게 확인할 수 있잖아요.

이런 좋은 기능과 역할을 하는데도 불구하고 unique index와 not null 만을 정의해서 primary key 기능을 대신하도록 할 이유는 없지 않을까요?
성능과 관련해서 말씀드리면, Not Null Constraint와 Unique 인덱스를 PK Constraint 대신 사용하는 것이 더 속도가 빠르다는 것은 전혀 근거없는 낭설에 불과합니다. 오히려 SQL 옵티마이저에게 더 많은 정보를 제공함으로써 더 좋은 실행계획을 만드는데 일조하게 되고 따라서 더 빨라지는 경우가 많겠지요...

# 토큰 기반의 인증

1. 사용자가 로그인함
2. 서버는 token에 sign하여 access_token을 전달함.
3. 사용자는 다른 api를 사용할 때 access_token을 포함하여 전달함.
4. 서버는 access_token을 token을 verify함
   
## JWT
JWT는 base64 방식으로 인코딩됨
3가지 파트로 되어있음 
- Header
  - signing algorithm
- Payload
  - 토큰에 대한 정보가 들어있음
  - id나 사용자 이름, 만료시간 등
- verify signature
  - 토큰 유효한지
  - 서버에서 어떤 값으로 사인하는지 256비트로 암호화되어있음  

### JWT의 서명 알고리즘들
#### symmetric(대칭) digital signature algorithm
- same secret key. is used to sign & verify
- for local use: internal services, where the secret key can be shared
- such as HS256, HS384, HS512
  - HS256 = HMAC + SHA256
  - HMAC = Hash-based Message Authentication Code

#### Asymmetric(비대칭) digital signature algorithm
- private key is ued to sign token
- public key is used to verify token
- for public use: internal service signs token, but external service needs to verify it
- RS256
  - RSA PKCSv1.5 + SHA256
    - PKCS(Public-Key Cryptography Standards)
- PS256
  - RSA PSS + SHA256
  - PSS is more secure than PKCS
    - PSS(Probabilistic Signature Scheme)
- ES256
  - ECDSA + SHA256
    - ECDSA(Elliptic Curve Digital Signature Algorithm)
    - 타원곡선암호 

## JWT의 문제
- 빈약한 서명 알고리즘
  - 개발자가 서명 알고리즘을 선택할 수 있음(rs256, ES256 등)
  - 예시로 RS256은 padding oracle attack을 허용할 수 있음
    - https://en.wikipedia.org/wiki/Padding_oracle_attack
  - ECDSA는 invalid-curve attack 공격 가능
  - 따라서 개발자가 보안적 깊은 지식이 없다면 어떤 알고리즘 선택이 가장 최선인지 알 수 없음
  - 이것이 문제가 됨. 너무 유연한 서명 알고리즘 선택권이 공격자의 기회가 될 수 있음
- 위조가 될 수 있음
  - 구현을 잘못하거나 라이브러리를 잘못 쓰면 공격자가 위조해서 뚫을 수 있음
    - 예를들면 JWT의 HEADER의 "alg"을 None으로 해버리면 서명 검증 단계를 통과해버림
    - 물론 이런 이슈들은 많은 라이브러리에서 수정되었음
-  potential attack
   - jwt의 알고리즘 해더를 symmetric(대칭) 알고리즘으로 사용함. (HS256 같은)
   - 서버가 만약 Asymmetric(비대칭) 알고리즘을 사용하면 문제가 발생함
    1. 서버가 RSA public 서명을 사용함
    2. 해커가 가짜 어드민 사용자의 토큰을 생성함
    3. 생성된 가짜 어드민 토큰의 헤더 알고리즘을 비대칭 알고리즘으로 설정함(hs256같은)
    5. 가짜 토큰을 public-key로 서명하여 서버에 전송
    6. 서버는 토큰을 검증할때 [3]에 의해 RSA알고리즘 대신에 HS256으로 인증해버림
    7. 이게 되는 이유는 public-key를 이미 해커가 토큰 payload에 서명해버려서 서명 검증 단계를 통과해버림
    8. 인증 권한 획득
   - 이를 예방하기 위해서는 서버 로직단에서 alg 헤더를 체크해야함. 

위 이유로 JWT는 잘못 설계된 토큰임

## PASETO(Platform-Agnostic Security Tokens)
JWT의 문제를 해결한 토큰
- Stronger algorithms
  - 개발자가 알고리즘을 선택할 필요가 없음 
  - PASETO 버전을 고르면 됨.
  - 로컬에서 대칭키를 사용한다면 encoding이 아니라 암호화를 진행함. 해커로부터 안전


### 개인적인 생각
새 프로젝트면 paseto 도입을 팀원과 고려해보면 좋을듯
그게 아니면 그냥 JWT 써도 큰 문제는 없을것

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
