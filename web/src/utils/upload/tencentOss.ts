import { FileUpload } from "@/types/idnex"
import { fileHash, fileSuffix } from "./fileManipulation"
import { uploadOssFile } from "@/apis/commonality"


export const tencentOssUpload = async (file: File, uploadConfig: FileUpload, dir: string, fragment?: boolean): Promise<any> => {
    return new Promise(async (resolve, reject) => {
            //直接上传
            // const name = await fileHash(file) + fileSuffix(file.name)
            const name = file.name
            const formData = new FormData()
            const key = `${name}`
            formData.append('interface', uploadConfig.interface)
            formData.append('name', name)
            formData.append('file', file)
            try {
                const response = await uploadOssFile(formData, uploadConfig)
                resolve({ path: response.data as string })
                console.log(response)
            } catch (err) {
                console.log(err)
                reject({ msg: '上传失败' })
            }
    })        
}