# gRPC

http/2 기반으로 전송되며 unary, bid streaming을 지원함

proto file 작성으로 protoc를 통해 클라이언트와 서버 둘 다 쓸 수 있는 RPC 생성됨
+ MSA에서 여러 언어들과 통신 가능
+ 높은 성능

# gRPC Gateway
protoc의 플러그인으로 Restful HTTP API를 gRPC로 변환해주어서 http request와 RPC 요청 모두를 한번에 serving 할 수 있다.

gRPC 호출을 HTTP JSON으로 변환하는데는 2가지가 있다.
+ in-process 변환
  + unary RPC에만 사용된다.
+ separate proxy server
  + unary와 streaming 둘 다 지원