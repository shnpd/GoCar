import { IAppOption } from "../../appoption"
import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"
import { rental } from "../../service/proto_gen/rental/rental_pb";
import { formatDuration, formatFee } from "../../utils/format";


interface Trip {
  id: string
  shortId: string
  start: string
  end: string
  duration: string
  fee: string
  distance: string
  status: string
}

interface MainItem {
  id: string
  navId: string
  navScrollId: string //实现左侧选择留白，即右侧选择30号元素，左侧把29号元素显示在第一个位置
  data: Trip
}

interface NavItem {
  id: string
  mainId: string
  label: string
}

interface MainItemQueryResult {
  id: string
  top: number
  dataset: {
    navId: string
    navScrollId: string
  }
}

const tripStatusMap = new Map([
  [rental.v1.TripStatus.IN_PROGRESS, '进行中'],
  [rental.v1.TripStatus.FINISHED, '已完成'],
])

Page({
  scrollStates: {
    mainItems: [] as MainItemQueryResult[]
  },

  layoutResolver: undefined as ((value: unknown) => void) | undefined,


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
    ],
    tripsHeight: 0,
    mainScroll: '',
    navItems: [] as NavItem[],
    mainItems: [] as MainItem[],
    navCount: 0,
    navSel: '',
    navScroll: '',
  },
  onLoad() {
    const layoutReady = new Promise((resolve) => {
      this.layoutResolver = resolve
    })

    // 等待layoutReady执行结束，layoutReady等待layoutResolver调用（等待layoutResolver调用）
    Promise.all([TripService.getTrips(), layoutReady]).then(([trips]) => {
      this.populateTrips(trips.trips!)
    })
    getApp<IAppOption>().globalData.userInfo.then(userInfo => {
      this.setData({
        avatarURL: userInfo.avatarUrl,
      })
    })
  },
  onReady() {
    wx.createSelectorQuery().select('#heading')
      .boundingClientRect(rect => {
        const height = wx.getSystemInfoSync().windowHeight - rect.height
        this.setData({
          tripsHeight: height,
          navCount: Math.round(height / 50)
        }, () => {
          if (this.layoutResolver) {
            this.layoutResolver('')
          }
        })
      }).exec()
  },
  populateTrips(trips: rental.v1.ITripEntity[]) {
    const mainItems: MainItem[] = []
    const navItems: NavItem[] = []
    let navSel = ''
    let preNav = ''
    for (let i = 0; i < trips.length; i++) {
      const trip = trips[i]
      const mainId = 'main-' + i
      const navId = 'nav-' + i
      const shortId = trip.id?.substr(trip.id.length - 6)
      if (!preNav) {
        preNav = navId
      }
      const tripData: Trip = {
        id: trip.id!,
        shortId: '****' + shortId,
        start: trip.trip?.start?.poiName || '未知',
        end: '',
        distance: '',
        duration: '',
        fee: '',
        status: tripStatusMap.get(trip.trip?.status!) || '未知',
      }

      const end = trip.trip?.end
      if (end) {
        tripData.end = end.poiName || '未知'
        tripData.distance = end.kmDriven?.toFixed(1) + '公里'
        tripData.fee = formatFee(end.feeCent || 0)
        const dur = formatDuration((end.timestampSec || 0) - (trip.trip?.start?.timestampSec || 0))
        tripData.duration = `${dur.hh}时${dur.mm}分`
      }
      mainItems.push({
        id: mainId,
        navId: navId,
        navScrollId: preNav,
        data: tripData
      })
      navItems.push({
        id: navId,
        mainId: mainId,
        label: shortId || '',
      })
      if (i === 0) {
        navSel = navId
      }
      preNav = navId
    }

    // display-multiple-items不能大于swiper-item数量，添加空白的swiper-item进行占位
    for (let i = 0; i < this.data.navCount - 1; i++) {
      navItems.push({
        id: '',
        mainId: '',
        label: '',
      })
    }
    this.setData({
      mainItems: mainItems,
      navItems: navItems,
      navSel
    }, () => {
      this.prepareScrollStates()
    })
  },
  prepareScrollStates() {
    wx.createSelectorQuery().selectAll('.main-item')
      .fields({
        id: true,
        dataset: true,
        rect: true,
      }).exec(res => {
        this.scrollStates.mainItems = res[0]
        console.log(this.scrollStates.mainItems)
      })
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
    // console.log(e)
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
  onNavItemTap(e: any) {
    console.log(e)
    const mainId: string = e.currentTarget?.dataset?.mainId
    const navId: string = e.currentTarget?.id
    if (mainId && navId) {
      this.setData({
        navSel: navId,
        mainScroll: mainId
      })
    }
  },
  onMainScroll(e: any) {
    const top: number = e.currentTarget?.offsetTop + e.detail?.scrollTop
    if (top === undefined) {
      return
    }
    // 获取当前右侧选择的元素
    const selItem = this.scrollStates.mainItems.find(
      v => v.top >= top
    )
    if (!selItem) {
      return
    }
    // 设定左边对应选择的元素
    this.setData({
      navSel: selItem.dataset.navId,
      // 同步滚动留空一元素
      navScroll: selItem.dataset.navScrollId
    })
  },
})