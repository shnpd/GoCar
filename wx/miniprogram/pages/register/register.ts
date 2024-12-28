import { ProfileService } from "../../service/profile"
import { rental } from "../../service/proto_gen/rental/rental_pb"
import { padString } from "../../utils/format"
import { routing } from "../../utils/routing"

function formatDate(millis: number) {
	const date = new Date(millis)
	const y = date.getFullYear()
	const m = date.getMonth() + 1
	const d = date.getDate()
	return `${padString(y)}-${padString(m)}-${padString(d)}`
}
Page({
	redirectURL: "",
	profileRefresher: 0,
	data: {
		state: rental.v1.IdentityStatus[rental.v1.IdentityStatus.UNSUBMITTED],
		genderIndex: 0,
		genders: ['未知', '男', '女'],
		licImgURL: '',
		birthDate: '1990-01-01',
		licNo: '',
		name: ''
	},
	renderProfile(p: rental.v1.IProfile) {
		this.setData({
			licNo: p.identity?.licNumber || '',
			name: p.identity?.name || '',
			genderIndex: p.identity?.gender || 0,
			birthDate: formatDate(p.identity?.birthDateMillis || 0),
			state: rental.v1.IdentityStatus[p.identityStatus || 0],
		})
	},
	onLoad(opt: Record<'redirect', string>) {
		const o: routing.registerOpts = opt
		if (o.redirect) {
			this.redirectURL = decodeURIComponent(o.redirect)
		}
		ProfileService.getProfile().then(p => this.renderProfile(p))
	},
	onUploadLic() {
		wx.chooseMedia({
			success: (res) => {
				if (res.tempFiles.length > 0) {
					this.setData({
						licImgURL: res.tempFiles[0].tempFilePath
					})
					const data = wx.getFileSystemManager().readFileSync(res.tempFiles[0].tempFilePath)
					wx.request({
						method:'PUT',
						url:'https://coolcar-1311261643.cos.ap-nanjing.myqcloud.com/account_1/676f88c84aea956a1e0ed05b?q-sign-algorithm=sha1&q-ak=AKIDaBWbPxHK7dvxiCQ4SZQ0JL6anslEWaPz&q-sign-time=1735362760%3B1735363760&q-key-time=1735362760%3B1735363760&q-header-list=host&q-url-param-list=&q-signature=36c40787d186b8e37a5e74ddc50b2aa52d5677ac',
						data,
						success:console.log,
						fail:console.error,
					})
				}
			}
		})
	},
	onGenderChange(e: any) {
		this.setData({
			genderIndex: parseInt(e.detail.value)
		})
	},
	onBirthdateChange(e: any) {
		this.setData({
			birthDate: e.detail.value
		})
	},
	onSubmit() {
		ProfileService.submitProfile({
			licNumber: this.data.licNo,
			name: this.data.name,
			gender: this.data.genderIndex,
			birthDateMillis: Date.parse(this.data.birthDate),
		}).then(p => {
			this.renderProfile(p)
			this.scheduleProfileRefresher()
		})
	},

	onUnload() {
		this.clearProfileRefresher()
	},

	scheduleProfileRefresher() {
		this.profileRefresher = setInterval(() => {
			ProfileService.getProfile().then(p => {
				this.renderProfile(p)
				if (p.identityStatus !== rental.v1.IdentityStatus.PENDING) {
					this.clearProfileRefresher()
				}
				if (p.identityStatus === rental.v1.IdentityStatus.VERIFIED) {
					this.onLicVerified()
				}
			})
		}, 1000);
	},
	clearProfileRefresher() {
		if (this.profileRefresher) {
			clearInterval(this.profileRefresher)
			this.profileRefresher = 0
		}
	},
	onResubmit() {
		ProfileService.clearProfile().then(p => this.renderProfile(p))
	},
	//TODO: 服务器审查通过
	onLicVerified() {
		if (this.redirectURL) {
			wx.redirectTo({
				url: this.redirectURL,
			})
		}
	}
})