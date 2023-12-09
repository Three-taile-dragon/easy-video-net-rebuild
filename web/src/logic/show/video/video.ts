import { danmakuApi, getVideoBarrageList, getVideoContributionByID, likeVideo, sendVideoBarrage } from "@/apis/contribution";
import globalScss from "@/assets/styles/global/export.module.scss";
import { useGlobalStore, useUserStore } from "@/store/main";
import { GetVideoBarrageListReq, GetVideoContributionByIDReq, LikeVideoReq, SendVideoBarrageReq, VideoInfo,BarrageListConvert } from "@/types/show/video/video";
import DPlayer, { DPlayerDanmakuItem, DPlayerVideoQuality } from "dplayer";
import Swal from 'sweetalert2';
import { Ref, UnwrapNestedRefs, reactive, ref } from "vue";
import { RouteLocationNormalizedLoaded, Router, useRoute, useRouter } from "vue-router";
import { numberOfViewers, responseBarrageNum } from './socketFun';
import TCPlayer from 'tcplayer.js';
import 'tcplayer.js/dist/tcplayer.min.css';

export const useVideoProp = () => {
  const route = useRoute()
  const router = useRouter()
  const global = useGlobalStore()
  const userStore = useUserStore()
  const videoRef = ref()
  const videoID = ref<number>(0)
  const videoInfo = reactive(<VideoInfo>{})
  const barrageListConvert = reactive(<BarrageListConvert>{})
  const barrageInput = ref("")
  const barrageListShow = ref(false)
  const videoBarrage = ref(true)
  const liveNumber = ref(0)
  //回复二级评论
  const replyCommentsDialog = reactive({
    show: false,
    commentsID: 0,
  })

  return {
    route,
    router,
    userStore,
    videoRef,
    videoID,
    videoInfo,
    barrageListConvert,
    barrageInput,
    barrageListShow,
    liveNumber,
    replyCommentsDialog,
    videoBarrage,
    global
  }
}

export const useSendBarrage = async (text: Ref<string>, dpaler: DPlayer, userStore: any, videoInfo: UnwrapNestedRefs<VideoInfo>, socket: WebSocket) => {
  const res = await sendVideoBarrage(<SendVideoBarrageReq>{
    author: userStore.userInfoData.username,
    color: 16777215,
    id: videoInfo.videoInfo.id.toString(),
    text: text.value,
    time: dpaler.video.currentTime,
    type: 0,
    token: userStore.userInfoData.token,
  })

  console.log(userStore.userInfoData)
  if (res.code != 0) {
    Swal.fire({
      title: "弹幕服务异常",
      heightAuto: false,
      confirmButtonColor: globalScss.colorButtonTheme,
      icon: "error",
    })
    return
  }
  const danmaku = <DPlayerDanmakuItem>{
    text: text.value,
    color: '#fff',
    type: 'right',
  };

  dpaler.danmaku.draw(danmaku);

  text.value = ""

  let data = JSON.stringify({
    type: "sendBarrage",
    data: ""
  })
  socket.send(data)

}

export const useLikeVideo = async (videoInfo: UnwrapNestedRefs<VideoInfo>) => {
  try {
    await likeVideo(<LikeVideoReq>{
      video_id: videoInfo.videoInfo.id
    })
    videoInfo.videoInfo.is_like = !videoInfo.videoInfo.is_like
  } catch (err) {
    Swal.fire({
      title: "点赞失败",
      heightAuto: false,
      confirmButtonColor: globalScss.colorButtonTheme,
      icon: "error",
    })
  }
}

export const useInit = async (videoRef: Ref, route: RouteLocationNormalizedLoaded, Router: Router, videoID: Ref<Number>, videoInfo: UnwrapNestedRefs<VideoInfo>, global: any) => {
  try {
    //绑定视频id
    if (!route.params.id) {
      Router.back()
      Swal.fire({
        title: "获取视频失败",
        heightAuto: false,
        confirmButtonColor: globalScss.colorButtonTheme,
        icon: "error",
      })
      Router.back()
      return
    }
    global.globalData.loading.loading = true
    videoID.value = Number(route.params.id)
    //得到视频信息
    const vinfo = await getVideoContributionByID(<GetVideoContributionByIDReq>{
      video_id: videoID.value
    })
    if (!vinfo.data) return false
    videoInfo.videoInfo = vinfo.data.videoInfo
    videoInfo.recommendList = vinfo.data.recommendList

    //得到清晰度列表
    let quality: DPlayerVideoQuality[] = []
    if (videoInfo.videoInfo.video) {
      quality = [...quality, {
        name: "1080P超清",
        url: videoInfo.videoInfo.video
      }]
    }
    if (videoInfo.videoInfo.video_720p) {
      quality = [...quality, {
        name: "720P高清",
        url: videoInfo.videoInfo.video_720p
      }]
    }
    if (videoInfo.videoInfo.video_480p) {
      quality = [...quality, {
        name: "480P标清",
        url: videoInfo.videoInfo.video_480p
      }]
    }
    if (videoInfo.videoInfo.video_360p) {
      quality = [...quality, {
        name: "360P流畅",
        url: videoInfo.videoInfo.video_360p
      }]
    }
    //获取视频弹幕信息
    const barrageList = await getVideoBarrageList(<GetVideoBarrageListReq>{
      id: videoID.value.toString()
    })
    if (!barrageList.data) return false
    videoInfo.barrageList = barrageList.data
    //获取当前用户信息
    const userStore = useUserStore()
    //初始化播放器
    const dp = new DPlayer({
      container: videoRef.value,
      loop: true, // 循环播放
      lang: "zh-cn", // 语言
      logo: "",
      autoplay: true,
      danmaku: {
        id: videoID.value.toString(),
        api: danmakuApi,
        token: userStore.userInfoData.token
      },
      mutex: false, // 互斥，阻止多个播放器同时播放
      video: {
        quality: quality,
        defaultQuality: 0,
        url: "不填", // 视频链接
        pic: videoInfo.videoInfo.cover
      },
    });
    global.globalData.loading.loading = false
    return dp
  } catch (err) {
    global.globalData.loading.loading = false
    console.log(err)
  }
}

export const tcPlayerInit = async (myRef: Ref, route: RouteLocationNormalizedLoaded, Router: Router, videoID: Ref<Number>, videoInfo: UnwrapNestedRefs<VideoInfo>, global: any) => {
  try {
    //绑定视频id
    if (!route.params.id) {
      Router.back()
      Swal.fire({
        title: "获取视频失败",
        heightAuto: false,
        confirmButtonColor: globalScss.colorButtonTheme,
        icon: "error",
      })
      Router.back()
      return
    }
    global.globalData.loading.loading = true
    videoID.value = Number(route.params.id)
    //得到视频信息
    const vinfo = await getVideoContributionByID(<GetVideoContributionByIDReq>{
      video_id: videoID.value
    })
    if (!vinfo.data) return false
    videoInfo.videoInfo = vinfo.data.videoInfo
    videoInfo.recommendList = vinfo.data.recommendList

    //获取视频弹幕信息
    const barrageList = await getVideoBarrageList(<GetVideoBarrageListReq>{
      id: videoID.value.toString()
    })
    if (!barrageList.data) return false
    videoInfo.barrageList = barrageList.data
    //获取当前用户信息
    const userStore = useUserStore()
    //初始化播放器
    var tp:any
    if (videoInfo.videoInfo.fileID != "") {
      tp = new TCPlayer(myRef.value, {
        fileID: videoInfo.videoInfo.fileID,
        appID: videoInfo.videoInfo.appID,
        psign: videoInfo.videoInfo.pSign,
        licenseUrl: videoInfo.videoInfo.licenseUrl,
        autoplay: true,
        plugins: {
          ProgressMarker: true,
          ContextMenu: {
            statistic: true
          },
          ContinuePlay: { // 开启续播功能
            // auto: true, //[可选] 是否在视频播放后自动续播
            text:'上次播放至 ', //[可选] 提示文案
            btnText: '恢复播放' //[可选] 按钮文案
          },
        }
      })
    }else{
      // 没有 file ID 即为本地资源
      tp = new TCPlayer(myRef.value, {
        autoplay: true,
        reportable: false,
        licenseUrl: videoInfo.videoInfo.licenseUrl,
        // sources: [{
        //   src: bestQualityUrl(videoInfo),
        //   type: 'video/mp4',
        // }],
        multiResolution:{
          // 配置多个清晰度地址
          sources:{
            'LOW':[{
              src: videoInfo.videoInfo.video_360p,
            }],
            'SD':[{
              src: videoInfo.videoInfo.video_480p,
            }],
            'HD':[{
              src: videoInfo.videoInfo.video_720p,
            }],
            'FHD':[{
              src: bestQualityUrl(videoInfo),
            }]
          },
          // 配置每个清晰度标签
          labels:{
            "LOW": '360P','SD':'480P','HD':'720P','FHD':'超清'
          },
          // 配置各清晰度在播放器组件上的顺序
          showOrder:['FHD','HD','SD',"LOW"],
          // 配置默认选中的清晰度
          defaultRes: 'FHD',
        },          
        plugins: {
          ProgressMarker: true,
          ContextMenu: {
            statistic: true
          }
        }
      })
    }
    console.log(tp)
    global.globalData.loading.loading = false
    return tp
  } catch (err) {
    global.globalData.loading.loading = false
    console.log(err)
  }
}

export const tcSendBarrage = async (text: Ref<string>, tcplayer: any,tcplayerBarrage: any, userStore: any, videoInfo: UnwrapNestedRefs<VideoInfo>, socket: WebSocket) => {
  const res = await sendVideoBarrage(<SendVideoBarrageReq>{
    author: userStore.userInfoData.username,
    color: 16777215,
    id: videoInfo.videoInfo.id.toString(),
    text: text.value,
    time: tcplayer.currentTime(),
    type: 0,
    token: userStore.userInfoData.token,
  })

  console.log(userStore.userInfoData)
  if (res.code != 0) {
    Swal.fire({
      title: "弹幕服务异常",
      heightAuto: false,
      confirmButtonColor: globalScss.colorButtonTheme,
      icon: "error",
    })
    return
  }
  var barrage = {
    "mode":1,
    "text": text.value,
    "size": 25,
    "color":'#ff0000',
  };
  // 即时发送弹幕
  tcplayerBarrage.send(barrage);

  text.value = ""
  let data = JSON.stringify({
    type: "sendBarrage",
    data: ""
  })
  socket.send(data)

}

export const barrageConvert = (videoInfo: UnwrapNestedRefs<VideoInfo>) =>{
  const barrageList = videoInfo.barrageList;
  const convertedList: BarrageListConvert[] = [];

  for (const barrage of barrageList) {
    const convertedBarrage: BarrageListConvert = {
      mode: 1, // 设置合适的弹幕模式
      text: barrage.text,
      stime: barrage.time *1000,
      size: 25, // 设置合适的字体大小
      color: '#ffffff' // 设置合适的字体颜色
    };

    convertedList.push(convertedBarrage);
  }

  return convertedList;
}

export const bestQualityUrl = (videoInfo: UnwrapNestedRefs<VideoInfo>) =>{
  if (videoInfo.videoInfo.video != ""){
    return videoInfo.videoInfo.video
  }else if (videoInfo.videoInfo.video_720p != "") {
    return videoInfo.videoInfo.video_720p
  }else if (videoInfo.videoInfo.video_480p != "") {
    return videoInfo.videoInfo.video_480p
  }else if (videoInfo.videoInfo.video_360p != "") {
    return videoInfo.videoInfo.video_360p
  }
  return ""
}

export const useWebSocket = (userStore: any, videoInfo: UnwrapNestedRefs<VideoInfo>, Router: Router, liveNumber: Ref<number>) => {
  let socket: WebSocket
  const open = () => {
    console.log("websocket 连接成功 ")
  }
  const error = () => {
    console.error("websocket 连接失败")
  }
  const getMessage = async (msg: any) => {
    let data = JSON.parse(msg.data)
    console.log(data)
    switch (data.type) {
      case "numberOfViewers":
        numberOfViewers(liveNumber, data.data.people)
        break;
      case "responseBarrageNum":
        responseBarrageNum(videoInfo)
        break;
    }
  }

  if (typeof (WebSocket) === "undefined") {
    Swal.fire({
      title: "您的浏览器不支持socket",
      heightAuto: false,
      confirmButtonColor: globalScss.colorButtonTheme,
      icon: "error",
    })
    Router.back()
    return
  } else {
    // 实例化socket
    socket = new WebSocket(import.meta.env.VITE_SOCKET_URL + "/ws/videoSocket?token=" + userStore.userInfoData.token + "&videoID=" + videoInfo.videoInfo.id)
    // 监听socket连接
    socket.onopen = open
    // 监听socket错误信息
    socket.onerror = error
    // 监听socket消息
    socket.onmessage = getMessage
  }

  return socket
}
