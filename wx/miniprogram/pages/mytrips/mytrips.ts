import { routing } from "../../utils/routing"

Page({
  data: {
    indicatorDots: true,
    autoPlay: false,
    interval: 3000,
    duration: 500,
    circular: true,
    multiItemCount: 1,
    prevMargin: '',
    nextMargin: '',
    vertical: false,
    current: 0,
    avatarURL: '',
    promotionItems: [
      {
        img: "https://img1.sycdn.imooc.com/65e82b7c00012aed17920764.jpg",
        promotionID: 1
      },
      {
        img: "https://img1.sycdn.imooc.com/65e532ed0001616e16000682.jpg",
        promotionID: 2
      },
      {
        img: "https://img1.sycdn.imooc.com/65e5367a0001bd8417920764.jpg",
        promotionID: 3
      },
      {
        img: "https://img1.sycdn.imooc.com/65dd5c790001d8dd17920764.jpg",
        promotionID: 4
      },
    ]
  },
  onShow() {
    const avatar = getApp<IAppOption>().globalData.avatarURL
    if (avatar) {
      this.setData({
        avatarURL: avatar
      })
    }
  },
  onSwiperChange(e: any) {
    // console.log(e)
  },
  onPromotionItemTap(e: any) {
    console.log(e)
  },
  onRegisterTap() {
    wx.navigateTo({
      url: routing.register()
    })
  },
  onBindChooseAvatar(e: any) {
    const avatar = e.detail.avatarUrl
    getApp<IAppOption>().globalData.avatarURL = avatar
    this.setData({ avatarURL: avatar, shareLocation: true })
  },
})