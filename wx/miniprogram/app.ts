import camelcaseKeys from "camelcase-keys"
import { IAppOption } from "./appoption"
import { auth } from "./service/proto_gen/auth/auth_pb"

// app.ts
App<IAppOption>({
  globalData: {},
  onLaunch() {
    // wx.request({
    //   url: 'http://localhost:8080/trip/trip123',
    //   method: 'GET',
    //   success: res => {
    //     const getTripRes = coolcar.GetTripResponse.fromObject(
    //       camelcaseKeys(res.data as object, {
    //         deep: true
    //       }))
    //     console.log(getTripRes)
    //     console.log('status is', coolcar.TripStatus[getTripRes.trip?.status!])
    //   },
    //   fail: console.error
    // })
    // 登录
    wx.login({
      success: res => {
        wx.request({
          url: 'http://localhost:8080/v1/auth/login',
          method: 'POST',
          data: {
            code: res.code,
          } as auth.v1.ILoginRequest,
          success: res => {
            const loginResp: auth.v1.ILoginResponse = auth.v1.LoginResponse.fromObject(camelcaseKeys(res.data as object))
            console.log(loginResp)
          },
          fail: console.error,
        })
        // 发送 res.code 到后台换取 openId, sessionKey, unionId
      },
    })
  },
})