export interface IAppOption {
  globalData: {
    avatarURL: string
    userInfo: Promise<WechatMiniprogram.UserInfo>
  }
  resolveUserInfo(userInfo: WechatMiniprogram.UserInfo): void
}