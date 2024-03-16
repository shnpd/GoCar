// components/modal/modal.ts
Component({

  /**
   * 组件的属性列表
   */
  properties: {
    showModal: Boolean,
    showCancel: Boolean,
    title: String,
    contents: String,
  },
  options: {
    addGlobalClass: true,
  },
  /**
   * 组件的初始数据
   */
  data: {

  },

  /**
   * 组件的方法列表
   */
  methods: {
    onCancel() {
      this.hideModal('cancel')
    },
    onOK() {
      console.log('111')
      this.hideModal('ok')
    },
    hideModal(res: 'ok' | 'cancel' | 'close') {
      this.setData({
        showModal: false
      })
      // 通知页面
      this.triggerEvent(res)
    }
  }
})