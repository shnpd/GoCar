import { ProfileService } from "../../service/profile"
import { rental } from "../../service/proto_gen/rental/rental_pb"
import { Coolcar } from "../../service/request"
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
		this.renderIdentity(p.identity!)
		this.setData({
			state: rental.v1.IdentityStatus[p.identityStatus || 0],
		})
	},

	renderIdentity(i?: rental.v1.IIdentity) {
		this.setData({
			licNo: i?.licNumber || '',
			name: i?.name || '',
			genderIndex: i?.gender || 0,
			birthDate: formatDate(i?.birthDateMillis || 0),
		})
	},
	onLoad(opt: Record<'redirect', string>) {
		const o: routing.registerOpts = opt
		if (o.redirect) {
			this.redirectURL = decodeURIComponent(o.redirect)
		}
		ProfileService.getProfile().then(p => this.renderProfile(p))
		ProfileService.getProfilePhoto().then(p => {
			this.setData({
				licImgURL: p.url || '',
			})
		})
	},
	onUploadLic() {
		wx.chooseMedia({
			success: async (res) => {
				if (res.tempFiles.length === 0) {
					return
				}
				this.setData({
					licImgURL: res.tempFiles[0].tempFilePath
				})
				// 获取图片上传地址
				const photoRes = await ProfileService.createProfilePhoto()
				// 上传图片
				if (!photoRes.uploadUrl) {
					return
				}
				await Coolcar.uploadFile({
					localPath: res.tempFiles[0].tempFilePath,
					url: photoRes.uploadUrl,
				})
				// 上传之后获取图片identity
				const identity = await ProfileService.completeProfilePhoto()
				this.renderIdentity(identity)
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
		ProfileService.clearProfilePhoto().then(() => {
			this.setData({
				licImgURL: '',
			})
		})
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