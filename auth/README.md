# 인증 서버 분리
Goal: 리소스 서버로부터 인증 서버를 분리

1. 리소스 서버로 회원가입/로그인 등 요청이 들어온다.
2. 리소스 서버는 auth server로 redirect
   - 클라이언트가 HTTP request를 리소스 서버가 받게되는데 리소스 서버는 grpc-gateway 역할을 하면서 auth server로 HTTP request를 리다이렉트
3. auth server가 grpc-gateway로 받은 요청을 처리하고 클라이언트에게 반환