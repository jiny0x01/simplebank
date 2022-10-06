# 인증 서버 분리
Goal: 리소스 서버로부터 인증 서버를 분리

1. 리소스 서버로 회원가입/로그인 등 요청이 들어온다.
2. 리소스 서버는 auth server로 redirect
   - 클라이언트가 HTTP request를 리소스 서버가 받게되는데 리소스 서버는 grpc-gateway 역할을 하면서 auth server로 HTTP request를 리다이렉트
3. auth server가 grpc-gateway로 받은 요청을 처리하고 클라이언트에게 반환

## Oauth flow

1. 클라이언트가 리소스 서버에 Oauth Authenticate 요청
2. 리소스 서버는 인증 서버로 리다이렉트
3. 인증 서버는 클라이언트에게 Authorization Request 요청
4. 클라이언트는 웹 페이지를 통해 로그인 인증하고 리소스 제공 범위(scope)를 설정하여 인가한다.
   1. 클라이언트는 인가 이후에 access_token과 refresh_token을 받는다.