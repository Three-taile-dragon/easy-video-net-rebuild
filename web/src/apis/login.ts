import httpRequest from "@/utils/requst"
import { loginReq , userInfoRes ,sendEmailReq,registReq ,forgetReq } from "@/types/login/login"


export const loginRequist = (params: loginReq) => {
    return httpRequest.post<userInfoRes>('/api/user/login', params);
}
export const regist = (params: registReq) => {
    return httpRequest.post<userInfoRes>('/api/user/register', params);
}
export const sendEmailVerificationCode = (params: sendEmailReq) => {
    return httpRequest.post('/api/user/getCaptcha', params);
}

export const sendEmailVerificationCodeByForget = (params: sendEmailReq) => {
    return httpRequest.post('/api/user/getCaptcha', params);
}

export const forgetRequist = (params: forgetReq) => {
    return httpRequest.post('/api/user/forget', params);
}

