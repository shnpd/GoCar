import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"

const centPerSec = 0.7
function formatDuration(sec: number): string {
  const padString = (n: number) => {
    return n < 10 ? '0' + n.toFixed(0) : n.toFixed(0)
  }
  let hour = Math.floor(sec / 3600)
  sec = sec - hour * 3600
  let min = Math.floor(sec / 60)
  sec = sec - min * 60
  return `${padString(hour)}:${padString(min)}:${padString(sec)}`
}
function formatFee(cents: number): string {
  return (cents / 100).toFixed(2)
}
Page({
  timer: undefined as number | undefined,
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
    console.log('current trip', o.trip_id)
    o.trip_id="67612a3189075d2ae73c9fd3"
    TripService.GetTrip(o.trip_id).then(console.log)
    this.setupTimer()
    // this.setupLocationUpdator()
  },
  onUnload() {
    if (this.timer) {
      clearInterval(this.timer)
    }
    // wx.stopLocationUpdate()
  },
  // setupLocationUpdator() {
  //   wx.startLocationUpdate({
  //     fail:console.error,
  //   })
  //   wx.onLocationChange((loc) => {
  //     this.setData({
  //       location: {
  //         latitude: loc.latitude,
  //         longitude: loc.longitude
  //       },
  //     })
  //   })
  // }
  setupTimer() {
    let elapsedSec = 0
    let cents = 0
    this.timer = setInterval(() => {
      elapsedSec++
      cents += centPerSec
      this.setData({
        elapsed: formatDuration(elapsedSec),
        fee: formatFee(cents)
      })
    }, 1000)
  },
  onEndTripTap() {
    wx.redirectTo({
      url: routing.mytrips(),
    })
  }
})