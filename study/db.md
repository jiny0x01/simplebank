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

# DBDOCS와 DBML(Database Markup Language)
dbdocs.io는 DB 스키마 정보를 웹상에서 시각화해주는 도구다.
https://dbdocs.io/docs

DBML은 오픈소스이며 DSL DB 스키마와 구조 정의하고 문서화하도록 설계되었다.
https://www.dbml.org/home/#what-can-i-do-now
CLI tool을 사용하여 관리 가능

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
   
