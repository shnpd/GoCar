export interface IAppOption {
    globalData: {
      userInfo: Promise<WechatMiniprogram.UserInfo>
    }
    userInfoReadyCallback?: WechatMiniprogram.GetUserInfoSuccessCallback,
  }