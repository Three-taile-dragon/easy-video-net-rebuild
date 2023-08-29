import { GetAttentionListReq, GetAttentionListRes, GetReleaseInformationReq, GetReleaseInformationRes, GetSpaceIndividualReq, GetSpaceIndividualRes, GetVermicelliListReq, GetVermicelliListRes } from "@/types/space/space";
import httpRequest from "@/utils/requst"
//获取个人空间信息
export const getSpaceIndividual = (params: GetSpaceIndividualReq) => {
    return httpRequest.post<GetSpaceIndividualRes>('/api/user/space/getSpaceIndividual', params);
}

//获取空间发布的视频和专栏
export const getReleaseInformation = (params: GetReleaseInformationReq) => {
    return httpRequest.post<GetReleaseInformationRes>('/api/user/space/getReleaseInformation', params);
}

//获取关注列表
export const getAttentionList = (params: GetAttentionListReq) => {
    return httpRequest.post<GetAttentionListRes>('/api/user/space/getAttentionList', params);
}

//获取粉丝列表
export const getVermicelliList = (params: GetVermicelliListReq) => {
    return httpRequest.post<GetVermicelliListRes>('/api/user/space/getVermicelliList', params);
}