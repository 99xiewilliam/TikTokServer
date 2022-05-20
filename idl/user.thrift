namespace go user

struct BaseResp {
    1:i64 status_code
    2:string status_message
    3:i64 service_time
}

struct CreateUserRequest {
    1:string user_name
    2:string password
    3:string name
}

struct CreateUserResponse {
    1:i64 user_id
    2:BaseResp base_resp
}

struct CheckUserRequest {
    1:string user_name
    2:string password
}

struct CheckUserResponse{
    1:i64 user_id
    2:BaseResp base_resp
}

struct CheckUserPresenceRequest {
    1:string user_name
}

struct CheckUserPresenceResponse {
    1:bool presence
    2:BaseResp base_resp
}

struct GetUserInfoRequest {
    1:i64 user_id
}

struct GetUserInfoResponse {
    1:i64 user_id
    2:string user_name
    3:string name
    4:BaseResp base_resp
}

service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req)
    CheckUserResponse CheckUser(1: CheckUserRequest req)
    CheckUserPresenceResponse CheckUserPresence(1: CheckUserPresenceRequest req)
    GetUserInfoResponse GetUserInfo(1: GetUserInfoRequest req)
}
