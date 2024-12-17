import { IAppOption } from "../../appoption"
import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"
import { rental } from "../../service/proto_gen/rental/rental_pb";


interface Trip {
  id: string
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

Page({
  scrollStates: {
    mainItems: [] as MainItemQueryResult[]
  },
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
  async onLoad() {
    const res = await TripService.GetTrips(rental.v1.TripStatus.FINISHED)
    this.populateTrips()
  },
  onReady() {
    wx.createSelectorQuery().select('#heading')
      .boundingClientRect(rect => {
        const height = wx.getSystemInfoSync().windowHeight - rect.height
        this.setData({
          tripsHeight: height,
          navCount: Math.round(height / 50)
        })
      }).exec()
  },
  populateTrips() {
    const mainItems: MainItem[] = []
    const navItems: NavItem[] = []
    let navSel = ''
    let preNav = ''
    for (let i = 0; i < 100; i++) {
      const mainId = 'main-' + i
      const navId = 'nav-' + i
      const tripId = (10001 + i).toString()
      if (!preNav) {
        preNav = navId
      }
      mainItems.push({
        id: mainId,
        navId: navId,
        navScrollId: preNav,
        data: {
          id: tripId,
          start: '东方明珠',
          end: '迪士尼',
          distance: '27.0公里',
          duration: '0时47分',
          fee: '128.00元',
          status: '已完成'
        }
      })
      navItems.push({
        id: navId,
        mainId: mainId,
        label: tripId
      })
      if (i === 0) {
        navSel = navId
      }
      preNav = navId
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