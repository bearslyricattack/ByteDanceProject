namespace go api

struct RegisterRequest {
	1: string username
	2: string password
	3: string name
}
struct RegisterResponse {
	1: string code
	2: string msg
	3: i64 userid
}

struct LoginRequest {
	1: string username
	2: string password

}
struct LoginResponse {
	1: string code
	2: string msg
	3: i64 userid
}

struct GetInfoRequest {
	1: i64 userid
}
struct GetInfoResponse {
	1: i64 userid
	2: string name
	3: string avatar
	4: string backgroundImage
	5: string signature
}

service User {
    RegisterResponse register(1: RegisterRequest req)
    LoginResponse login(1: LoginRequest req)
    GetInfoResponse get(1:GetInfoRequest req)
}
