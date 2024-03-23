export interface IAppOption {
    globalData: {
      userInfo?: WechatMiniprogram.UserInfo,
      avatarURL?:string,
    }
    userInfoReadyCallback?: WechatMiniprogram.GetUserInfoSuccessCallback,
  }