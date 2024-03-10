import { routing } from "../../utils/routing"

Page({
  redirectURL: "",
  data: {
    state: 'UNSUBMITTED' as 'UNSUBMITTED' | 'PENDING' | 'VERIFIED',
    genderIndex: 0,
    genders: ['未知', '男', '女'],
    licImgURL: '',
    birthDate: '1990-01-01',
    licNo: '',
    name: ''
  },
  onLoad(opt: Record<'redirect', string>) {
    const o: routing.registerOpts = opt
    if (o.redirect) {
      this.redirectURL = decodeURIComponent(o.redirect)
    }
  },
  onUploadLic() {
    wx.chooseMedia({
      success: (res) => {
        if (res.tempFiles.length > 0) {
          this.setData({
            licImgURL: res.tempFiles[0].tempFilePath
          })
          //TODO: upload image,假设一秒后服务器返回了驾照数据
          setTimeout(() => {
            this.setData({
              licNo: '123456',
              name: 'ssss',
              genderIndex: 1,
              birthDate: '2000-01-01'
            })
          }, 1000)
        }
      }
    })
  },
  onGenderChange(e: any) {
    this.setData({
      genderIndex: e.detail.value
    })
  },
  onBirthdateChange(e: any) {
    this.setData({
      birthDate: e.detail.value
    })
  },
  onSubmit() {
    this.setData({
      state: 'PENDING'
    })
    //TODO：模拟服务器审查
    setTimeout(() => {
      this.onLicVerified()
    }, 3000);
  },
  onResubmit() {
    this.setData({
      state: 'UNSUBMITTED',
      licImgURL: '',
    })
  },
  //TODO: 服务器审查通过
  onLicVerified() {
    this.setData({
      state: 'VERIFIED',
    })
    if (this.redirectURL) {
      wx.redirectTo({
        url: this.redirectURL,
      })
    }
  }
})