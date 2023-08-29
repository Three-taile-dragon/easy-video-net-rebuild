export interface loginReq  {
    name :string,
    password : string
}

export interface userInfoRes {
    username: string
    id : number
    photo : string
    token : string
    name : string
    created_at : string
}

export interface registReq {
    name :string
    password :string
    password2: string
    captcha :string
    email :string
}

export interface forgetReq {
    password :string
    password2: string
    captcha :string
    email :string
}

export interface sendEmailReq{
    email : string
}

export type sendEmailType = "regist" |"modify" 


export interface sendEmailInfo{
    btnText :string
    isPleaseClick : Boolean
    
}