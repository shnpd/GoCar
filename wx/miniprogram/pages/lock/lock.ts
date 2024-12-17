import { IAppOption } from "../../appoption"
import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"

const shareLocationKey = "share_location"

Page({
  carID:"",
  data: {
    shareLocation: false,
    avatarURL: "",
  },
  onLoad(opt: Record<'car_id', string>) {
    const o: routing.lockOpts = opt
    this.carID = o.car_id
    this.setData({
      shareLocation: wx.getStorageSync(shareLocationKey) || false,
    })
  },
  onBindChooseAvatar(e: any) {
    const avatar = e.detail.avatarUrl
    getApp<IAppOption>().globalData.avatarURL = avatar
    this.setData({ avatarURL: avatar, shareLocation: true })
  },
  onShareLocation(e: any) {
    const shareLocation: boolean = e.detail.value
    wx.setStorageSync(shareLocationKey, shareLocation)
    getApp<IAppOption>().globalData.avatarURL = shareLocation ? this.data.avatarURL : ""
  },
  onUnlockTap() {
    wx.getFuzzyLocation({
      type: "gcj02",
      success: async loc => {
        console.log('starting a trip', {
          location: {
            latitude: loc.latitude,
            longitude: loc.longitude
          },
          //TODO：需要双向绑定，这里即使选择不展示头像仍然会展示头像
          avatarURL: this.data.shareLocation ? this.data.avatarURL : "",
        })
        if (!this.carID) {
          console.error("carID is empty")
          return
        }
        const trip = await TripService.CreateTrip({
          start:loc,
          carId:this.carID,
        })
        if (!trip.id){
          console.error("trip id is empty")
          return
        }
        wx.showLoading({
          title: "开锁中",
          mask: true
        })
        setTimeout(() => {
          wx.redirectTo({
            url: routing.driving({
              trip_id: trip.id,
            }),
            complete: () => {
              wx.hideLoading()
            }
          })
        }, 2000);
      },
      fail: res => {
        console.log(res.errMsg)
        wx.showToast({
          icon: "none",
          title: "请前往设置页授权位置信息"
        })
      }
    })
  }
})