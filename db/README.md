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
