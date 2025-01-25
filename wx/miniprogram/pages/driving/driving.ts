import { rental } from "../../service/proto_gen/rental/rental_pb"
import { TripService } from "../../service/trip"
import { formatDuration, formatFee } from "../../utils/format"
import { routing } from "../../utils/routing"

const updateIntervalSec = 5


function durationStr(sec:number){
  const dur = formatDuration(sec)
  return `${dur.hh}:${dur.mm}:${dur.ss}`
}

Page({
  timer: undefined as number | undefined,
  tripID: "",
  data: {
    location: {
      latitude: 32.92,
      longitude: 114.32
    },
    elapsed: "00:00:00",
    fee: "0.00"
  },
  onLoad(opt: Record<'trip_id', string>) {
    const o: routing.drivingOpts = opt
    this.tripID = o.trip_id
    this.setupLocationUpdator()
    this.setupTimer(o.trip_id)
  },
  onUnload() {
    if (this.timer) {
      clearInterval(this.timer)
    }
    // wx.stopLocationUpdate()
  },
  setupLocationUpdator() {
    wx.startLocationUpdate({
      fail: console.error,
    })
    wx.onLocationChange((loc) => {
      this.setData({
        location: {
          latitude: loc.latitude,
          longitude: loc.longitude
        },
      })
    })
  },
  async setupTimer(tripID: string) {
    // 初始更新
    // 更新最新状态
    const trip = await TripService.getTrip(tripID)
    if(trip.status !== rental.v1.TripStatus.IN_PROGRESS){
      console.error("trip is not in progress")
      return
    }
    // 自从上次更新经过了多少秒
    let secSinceLastUpdate = 0
    // 上次更新时行程已经过了多少秒
    let lastUpdateDurationSec = trip.current!.timestampSec! - trip.start!.timestampSec!
    this.setData({
      elapsed: durationStr(lastUpdateDurationSec),
      fee: formatFee(trip.current!.feeCent!),
    })


    // 每5s更新一次
    this.timer = setInterval(() => {
      secSinceLastUpdate++
      if (secSinceLastUpdate % 5 === 0) {
        TripService.getTrip(tripID).then(trip => {
          lastUpdateDurationSec = trip.current!.timestampSec! - trip.start!.timestampSec!
          secSinceLastUpdate = 0
          this.setData({
            fee: formatFee(trip.current!.feeCent!),
          })
        }).catch(console.error)
      }
      this.setData({
        elapsed: durationStr(lastUpdateDurationSec + secSinceLastUpdate),
      })
    }, 1000)
  },
  onEndTripTap() {
    TripService.finishTrip(this.tripID).then(() => {
      wx.redirectTo({
        url: routing.mytrips(),
      })
    }).catch(err => {
      console.error(err)
      wx.showToast({
        title: "结束行程失败",
        icon: "none",
      })
    })
  }
})