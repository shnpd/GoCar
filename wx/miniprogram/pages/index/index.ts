import { IAppOption } from "../../appoption"
import { CarService } from "../../service/car"
import { ProfileService } from "../../service/profile"
import { rental } from "../../service/proto_gen/rental/rental_pb"
import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"

// index.ts
Page({

    isPageShowing: false,
    socket: undefined as WechatMiniprogram.SocketTask | undefined,
    data: {
        avatarURL: "",
        setting: {
            skew: 0,
            rotate: 0,
            showLocation: true,
            showScale: true,
            subKey: '',
            layerStyle: -1,
            enableZoom: true,
            enableScroll: true,
            enableRotate: false,
            showCompass: false,
            enable3D: false,
            enableOverlooking: false,
            enableSatellite: false,
            enableTraffic: false,
        },
        location: {
            latitude: 23.099994,
            longitude: 113.324520,
        },
        scale: 10,
        markers: [
            {
                iconPath: "/resources/car.png",
                id: 0,
                latitude: 23.099994,
                longitude: 113.324520,
                width: 50,
                height: 50
            },
            {
                iconPath: "/resources/car.png",
                id: 1,
                latitude: 23.099994,
                longitude: 114.324520,
                width: 50,
                height: 50
            },
        ],
    },
    async onLoad() {
        // 创建websocket连接
        let msgReceived = 0
        this.socket = CarService.subscribe((msg) => {
            msgReceived++
            console.log(msg)
        })

        const userInfo = await getApp<IAppOption>().globalData.userInfo
        this.setData({ avatarURL: userInfo.avatarUrl })
    },
    onMyLocationTap() {
        this.socket?.close({})
        wx.getFuzzyLocation({
            type: 'gcj02',
            success: (res: any) => {
                console.log(res.errMsg)
                this.setData({
                    location: {
                        latitude: res.latitude,
                        longitude: res.longitude
                    },
                })
            },
            fail: (res: any) => {
                console.log(res.errMsg)
                wx.showToast({
                    icon: "none",
                    title: "请前往设置页授权"
                })
            }
        })
    },

    async onShow() {
        const userInfo = await getApp<IAppOption>().globalData.userInfo
        this.setData({ avatarURL: userInfo.avatarUrl })
    },
    onHide() {

    },
    async onScanTap() {
        // 在有进行中的行程时不能扫码直接跳转到行程页面
        const trips = await TripService.getTrips(rental.v1.TripStatus.IN_PROGRESS)
        if (trips.trips?.length || 0 > 0) {
            await this.selectComponent('#tripModal').showModal()
            wx.navigateTo({
                url: routing.driving({
                    trip_id: trips.trips![0].id!,
                })
            })
            return
        }
        wx.scanCode({
            success: async () => {
                const carID = '6778d498ef7c52b63b87a130'
                const lockURL = routing.lock({
                    car_id: carID
                })
                const prof = await ProfileService.getProfile()
                if (prof.identityStatus === rental.v1.IdentityStatus.VERIFIED) {
                    wx.navigateTo({
                        url: lockURL,
                    })
                } else {
                    await this.selectComponent('#licModal').showModal()
                    //TODO:get car id from scan result
                    wx.navigateTo({
                        url: routing.register({
                            redirectURL: lockURL
                        })
                    })
                }
            }
        })
    },
    moveCars() {
        //获取map对象
        const map = wx.createMapContext("map")
        const dest = {
            latitude: 23.099994,
            longitude: 113.324520,
        }
        const moveCar = () => {
            dest.latitude += 0.1
            dest.longitude += 0.1
            map.translateMarker({
                destination: {
                    latitude: dest.latitude,
                    longitude: dest.longitude,
                },
                markerId: 0,
                autoRotate: false,
                rotate: 0,
                duration: 5000,
                animationEnd: () => {
                    if (this.isPageShowing)
                        moveCar()
                }
            })
        }
        moveCar()
    },
    onMyTripsTap() {
        wx.navigateTo({
            url: routing.mytrips()
        })
    },
    onModalOK() {
        console.log('ok clicked')
    }

})
