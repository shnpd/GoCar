import { IAppOption } from "../../appoption"
import { CarService } from "../../service/car"
import { car } from "../../service/proto_gen/car/car_pb"
import { rental } from "../../service/proto_gen/rental/rental_pb"
import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"

const shareLocationKey = "share_location"

Page({
  carID: "",
  carRefresher: 0,
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
    this.data.shareLocation = e.detail.value
    wx.setStorageSync(shareLocationKey, this.data.shareLocation)
  },
  onUnlockTap() {
    wx.getFuzzyLocation({
      type: "gcj02",
      success: async loc => {
        if (!this.carID) {
          console.error("carID is empty")
          return
        }
        let trip: rental.v1.ITripEntity
        try {
          trip = await TripService.createTrip({
            start: loc,
            carId: this.carID,
            avatarUrl: this.data.shareLocation ? this.data.avatarURL : "",
          })
          if (!trip.id) {
            console.error("trip id is empty")
            return
          }
        } catch (err) {
          wx.showToast({
            title: "创建行程失败",
            icon: "none",
          })
          return
        }

        wx.showLoading({
          title: "开锁中",
          mask: true
        })

        // 轮询汽车状态，如果开锁成功则跳转到驾驶页
        this.carRefresher = setInterval(async () => {
          const c = await CarService.getCar(this.carID)
          if (c.status === car.v1.CarStatus.UNLOCKED) {
            this.clearCarRefresher()
            wx.redirectTo({
              url: routing.driving({
                trip_id: trip.id!,
              }),
              complete: () => {
                wx.hideLoading()
              }
            })
          }
        }, 2000)
      },
      fail: (res: any) => {
        console.log(res.errMsg)
        wx.showToast({
          icon: "none",
          title: "请前往设置页授权位置信息"
        })
      }
    })
  },

  onUnload(){
    this.clearCarRefresher()
    wx.hideLoading()
  },
  clearCarRefresher() {
    if (this.carRefresher){
      clearInterval(this.carRefresher)
      this.carRefresher = 0
    }
  },

})