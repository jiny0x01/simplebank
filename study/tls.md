
## SSL/TLS
Netscape에서 SSL을 만들었는데 지금은 다 deprecated 되었음
TLS은 SSL의 업그레이드 버전으로 현재는 2018년에 배포된 1.3버전이 쓰임
+ website: HTTPS = HTTP + TLS
+ email: SMTPS = SMTP + TLS
+ File transfer: FTPS = FTP + TLS
  
## TLS를 사용하는 이유

+ Authentication(인증)
  + 상호간 신원 검증을 비대칭 암호로 인증지
+ Confidentiality(기밀성)
  + 대칭 암호화를 사용하여 비인증된 접근으로부터 교환된 데이터를 보호합니다.
+ Integrity(완전성)
  + 메시지 인증코드로 통신도중에 데이터 변조를 방지함

## 작동방식
2page로 작동한다.

1. Handshake protocol (인증)
  + 클라이언트-서버 TLS 버전 선택(1.2 or 1.3)
  + 암호화 알고리즘 선택
  + 비대칭 암호로 인증
  + 대칭암호화(기밀성)에 쓸 공용 secret key 설정
2. Record protocol(기밀성, 완정성)
  + 전송되는 모든 메시지는 handshake 과정에서 설정한 secret key로 암호화된다.
  + 암호화된 메시지는 상대방에게 전송되며 전송 중 위변조 여부를 알기 위해 확인된다.
  + 위변조가 안되었다면 설정한 대칭 secret key로 복호화한다.

대칭 암호와 비대칭 암호를 같이 사용하는 이유
+ 대칭 암호는 인증에 사용할 수 없고 key 공유가 어렵다
+ 비대칭 암호는 대칭 암호보다(100~10000배) 느리고 bulk encryption(Record protocol 단계)에 부적합하다.

## BIT FLIPPING ATTACK
1. Alice가 은행에 100달러를 누군가에게 보내려한다.
2. 이를 은행과 Alice가 갖고 있는 대칭 secret key로 암호화하여 전달한다.
3. 이때 해커가 암호화된 메시지를 가로챈다.
4. 암호화된 데이터의 0은 1로, 1은 0으로 바꾼다.
5. 위조된 데이터를 은행이 받고 복호화 하면 100달러(원본 데이터)가 900달러(위조데이터)로 바뀌어 있을 수 있다.

## Authenticated Encryption으로 bit flipping 방지
암호화된 메시지를 인증하는 방법

### 암호화 과정
1. 평문메시지를 대칭 암호화 알고리즘(aes-256, chacha20)으로 암호화. 암호화 알고리즘도 공유된 secret key를 가짐
2. random nonce나 initialization vector(IV)를 암호화 알고리즘의 input으로 준다.
3. 평문 메시지는 암호화 알고리즘에 의해 암호화된 메시지가 된다.
4. 암호화된 메시지, 비밀키, nonce는 MAC 알고리즘(GMAC{AES256 쓸 경우}, POLY1305{chacha20 쓸 경우})의 입력값이 된다.
  + Address, Ports, Protocol version, Sequence number와 같은 데이터는 클라이언트 서버 양측이 알고 있는 데이터이므로 이 데이터를 MAC algorithm의 input으로 사용하기도 한다.
5. MAC 알고리즘은 암호화된 메시지가 해싱되고 이는 메시지 인증 코드가 된다.
6. encrypt 메시지에 MAC 알고리즘으로 생성된 메시지 인증 코드를 붙여서 bob에게 보낸다.
  
### 복호화 과정

1. MAC 태그가 붙은 암호화된 메시지를 MAC 알고리즘의 태그를 구해서 tag를 떼어낸다.
2. 사전에 합의한 MAC 알고리즘에 input을 넣고 output을 떼어낸 tag와 비교한다. 값이 다르면 변조된 것이다.
3. 공유된 비밀키와 nonce값으로 대칭암호화된 메시지를 복호화 한다.

### 비밀키는 어떻게 공유하는가?
비대칭키(public-key)를 이용하여 공유하면 된다.
+ Diffie-Hellman Ephemeral(DHE)
+ Elliptic Curve Diffe-Hellman Ephemeral(ECDHE)

## Certificate Signing
1. Bob이 public-private key를 생성하고 CA에 certificate signing을 자신의 신원과 함께 요청함.
2. CA는 Bob의 신원을 증명하고 CA의 private key로 bob의 certificate에 sign함
3. Bob이 Alice에게 CA가 sign한 자신의 public key가 담긴 certificate를 보냄
4. Alice는 CA의 public key로 verify함.