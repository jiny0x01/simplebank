
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