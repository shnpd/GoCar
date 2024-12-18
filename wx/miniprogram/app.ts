
import { IAppOption } from "./appoption"

import { Coolcar } from "./service/request"

// app.ts
App<IAppOption>({ 
  globalData: {},
  async onLaunch() {
    // 登录
    Coolcar.login()
    
  },
})