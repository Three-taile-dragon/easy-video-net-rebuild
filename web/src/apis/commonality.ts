import { GetFullPathOfImageRrq, GetuploadingDirReq, GetuploadingDirRes, GetUploadingMethodReq, GetUploadingMethodRes, GteossStsres, SearchReq, SearchRes, UploadCheckReq, UploadCheckRes, UploadMergeReq, UploadMergeRes } from "@/types/commonality/commonality";
import { FileSliceUpload, FileUpload } from "@/types/idnex";
import httpRequest from "@/utils/requst";

//aliyun 获取oss sts 信息
export const gteossSTS = () => {
    return httpRequest.post<GteossStsres>('/commonality/ossSTS');
}

//获取上传方法
export function getuploadingMethod(params: GetUploadingMethodReq) {
    return httpRequest.post<GetUploadingMethodRes>('/api/commonality/uploadingMethod', params);
}

//获取保存地址
export function getuploadingDir(params: GetuploadingDirReq) {
    return httpRequest.post<GetuploadingDirRes>('/api/commonality/uploadingDir', params);
}

//获取图片完整路径
export function getFullPathOfImage(params: GetFullPathOfImageRrq) {
    return httpRequest.post<string>('/api/commonality/getFullPathOfImage', params);
}

//搜索功能
export function search(params: SearchReq) {
    return httpRequest.post<SearchRes>('/api/commonality/search', params);
}

//上传文件
export const uploadFile = (params: any, uploadConfig: FileUpload) => {
    return httpRequest.upload('/api/commonality/upload', params, uploadConfig);
}

//分片上传文件
export const UploadSliceFile = (params: any, uploadConfig: FileSliceUpload) => {
    return httpRequest.uploadSlice('/api/commonality/UploadSlice', params, uploadConfig);
}

//上传文件验证(验证操作)
export const uploadCheck = (params: UploadCheckReq) => {
    return httpRequest.post<UploadCheckRes>('/api/commonality/uploadCheck', params);
}

//上传文件验证(合并操作)
export const uploadMerge = (params: UploadMergeReq) => {
    return httpRequest.post<UploadMergeRes>('/api/commonality/uploadMerge', params);
}

// //注册媒体资源
// export const registerMedia = (params: RegisterMediaReq) => {
//     return httpRequest.post<RegisterMediaRes>('/commonality/registerMedia', params);
// }

//上传文件
export const uploadOssFile = (params: any, uploadConfig: FileUpload) => {
    return httpRequest.upload('/api/commonality/uploadOss', params, uploadConfig);
}