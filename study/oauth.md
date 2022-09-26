# Oauth 2.0 Implement
+ 목표
	+ https://www.rfc-editor.org/rfc/rfc6749 oauth2.0 framework spec 문서를 보고 client 구현하여 provider와 통신
	+ resource server와 authorization server 운영
+ 요구사항
	+ google oauth 인증으로 simple-bank 계정 생성 구현
	+ scope
		+ 기본
			+ 일반 사용자 권한 부여
		+ 전체권한부여
			+ admin과 동일한 권한 부여

## 분석
2022.09.11(03:40 pm)
### Oauth를 사용하게 된 배경
현대의 client-server 인증 모델은 server의 인증을 통해 client는 제한된 리소스를 사용합니다.
third-party앱이 제한된 리소스에 접근을 제공하기 위해 리소스 소유자는 third party에게 credentials(자격증명)을 공유합니다.
여기에 몇 가지 문제와 한계가 있다.
+ third-party 앱을 나중에 사용할 수 있도록 리소스 소유자의 자격증명(일반적으로 평문 text 비밀번호)를 저장해야한다.
+ 서버는 password에 내제된 보안 취약점에도 불구하고 비밀번호 인증을 지원해야한다.
+ third-party 앱은 리소스 소유자의 보호된 리소스에 광범위하게 엑세스 할 수 있고, 리소스 소유자가 기간을 제한하거나 제한된 하위 리소스들에 접근할도록만 할 수 없다.
+ 리소스 소유자는 각각의 third-party에 모든 third-party에게 access를 제한하지 않고서는 access를 철회할 수 없으며 third-party의 pasword를 바꾸는 방법으로 접근을 취소해야한다.
+ third-party중 하나라도 비밀번호가 뚫리면 해당 비밀번호를 사용하는 모든 데이터를 보호할 수 없다.

Oauth는 위 내용들을 바탕으로 authorization layer을 도입하여 클라이언트의 역할과 리소스 소유자의 역할을 구분한다.

### 역할
+ resource owner
	+ 제한된 리소스에 접근할 수 있는 존재. 사람이면 사용자에 해당
+ resource server
	+ 제한된 리소스를 호스팅하는 서버. access token을 사용하여 제한된 리소스를 수락하고 응답할 수 있다.
+ client
	+ 리소스 소유자 및 해당 권한을 대신하여 보호된 리소스를 요청을 수행하는 주체. 
	+ 서버, 데스크탑, 다른 디바이스가 클라이언트가 될 수 있다.
+ authorization server
	+ 서버가 리소스 소유자를 인증하고 인증이 성공되면 client에게 access token을 발급한다.
+

[1.2](https://www.rfc-editor.org/rfc/rfc6749#section-1.2).  Protocol Flow

```text
explicit flow 

     +--------+                               +---------------+
     |        |--(A)- Authorization Request ->|   Resource    |
     |        |                               |     Owner     |
     |        |<-(B)-- Authorization Grant ---|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(C)-- Authorization Grant -->| Authorization |
     | Client |                               |     Server    |
     |        |<-(D)----- Access Token -------|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(E)----- Access Token ------>|    Resource   |
     |        |                               |     Server    |
     |        |<-(F)--- Protected Resource ---|               |
     +--------+                               +---------------+

                     Figure 1: Abstract Protocol Flow

				 (a): resource owner로 부터 코드 달라고 요청
				 (b): 코드 받음()
				 (c): 받은 코드로 authorization server에 access token 달라고함
				 (d): access token 받음
				 (e): access token으로 resource 요청
				 (f): 유효한 access token이면 제한된 리소스(authroization) 허용
```

```text
  
  +--------+                                           +---------------+
  |        |--(A)------- Authorization Grant --------->|               |
  |        |                                           |               |
  |        |<-(B)----------- Access Token -------------|               |
  |        |               & Refresh Token             |               |
  |        |                                           |               |
  |        |                            +----------+   |               |
  |        |--(C)---- Access Token ---->|          |   |               |
  |        |                            |          |   |               |
  |        |<-(D)- Protected Resource --| Resource |   | Authorization |
  | Client |                            |  Server  |   |     Server    |
  |        |--(E)---- Access Token ---->|          |   |               |
  |        |                            |          |   |               |
  |        |<-(F)- Invalid Token Error -|          |   |               |
  |        |                            +----------+   |               |
  |        |                                           |               |
  |        |--(G)----------- Refresh Token ----------->|               |
  |        |                                           |               |
  |        |<-(H)----------- Access Token -------------|               |
  +--------+           & Optional Refresh Token        +---------------+

               Figure 2: Refreshing an Expired Access Token
```


## 구현

OAuth 회원가입
+ email
+ Fullname
+ Access_token
+ Refresh_Token
+ provider(google, kakao, etc)