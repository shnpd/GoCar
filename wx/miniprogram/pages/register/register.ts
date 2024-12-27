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

	onUnload(){
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