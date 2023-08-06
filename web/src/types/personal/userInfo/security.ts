export interface changePasswordReq {
    captcha : string
    password : string
    confirm_password : string
}

export interface sendEmailInfo{
    btnText :string
    isPleaseClick : Boolean
    
}