import camelcaseKeys from "camelcase-keys"
import { auth } from "./proto_gen/auth/auth_pb"

export namespace Coolcar {
    export const serverAddr = 'http://localhost:8080'
    export const wsAddr = 'ws://localhost:9090'
    const AUTH_ERR = 'AUTH_ERR'

    const authData = {
        token: '',
        expiryMs: 0,
    }

    export interface RequestOption<REQ, RES> {
        method: 'GET' | 'PUT' | 'POST' | 'DELETE'
        path: string
        data?: REQ
        respMarshaller: (r: object) => RES
    }

    export interface AuthOption {
        attachAuthHeader: boolean
        retryOnAuthError: boolean
    }

    export async function sendRequestWithAuthRetry<REQ, RES>(o: RequestOption<REQ, RES>, a?: AuthOption): Promise<RES> {
        const authOpt = a || { attachAuthHeader: true, retryOnAuthError: true }
        try {
            await login()
            return sendRequest(o, authOpt)
        } catch (err) {
            // 第二次登录失败则直接抛出错误
            if (err === AUTH_ERR && authOpt.retryOnAuthError) {
                authData.token = ''
                authData.expiryMs = 0
                return sendRequestWithAuthRetry(o, {
                    attachAuthHeader: authOpt.attachAuthHeader,
                    retryOnAuthError: false,
                })
            } else {
                throw err
            }
        }
    }

    export async function login() {
        // 如果token有效则无需再次登录直接返回
        if (authData.token && authData.expiryMs >= Date.now()) {
            return
        }
        // 如果token无效则重新发送登录请求
        const wxResp = await wxLogin()
        const reqTimeMs = Date.now()
        const resp = await sendRequest<auth.v1.ILoginRequest, auth.v1.ILoginResponse>({
            method: 'POST',
            path: '/v1/auth/login',
            data: {
                code: wxResp.code,
            },
            respMarshaller: auth.v1.LoginResponse.fromObject,
        }, {
            attachAuthHeader: false,
            retryOnAuthError: false,
        })
        authData.token = resp.accessToken!
        authData.expiryMs = reqTimeMs + resp.expiresIn! * 1000
    }

    function sendRequest<REQ, RES>(o: RequestOption<REQ, RES>, a: AuthOption): Promise<RES> {
        return new Promise((resolve, reject) => {
            const header: Record<string, any> = {}
            // 用户选择添加token
            if (a.attachAuthHeader) {
                // 如果token有效则加上token
                if (authData.token && authData.expiryMs >= Date.now()) {
                    header.authorization = 'Bearer ' + authData.token
                } else {
                    // token无效则返回错误
                    reject(AUTH_ERR)
                    return
                }
            }
            header['Grpc-Metadata-impersonate-account-id']="bad_account"

            wx.request({
                url: serverAddr + o.path,
                method: o.method,
                data: o.data,
                header,
                success: res => {
                    // 如果返回401则token无效
                    if (res.statusCode === 401) {
                        reject(AUTH_ERR)
                    } else if (res.statusCode >= 400) {
                        // 返回其他错误码则直接返回
                        reject(res)
                    } else {
                        // 正常返回则解析数据
                        resolve(o.respMarshaller(
                            camelcaseKeys(res.data as object, {
                                deep: true,
                            })))
                    }

                },
                fail: reject,
            })
        })
    }

    function wxLogin(): Promise<WechatMiniprogram.LoginSuccessCallbackResult> {
        return new Promise((resolve, reject) => {
            wx.login({
                success: resolve,
                fail: reject
            })
        })
    }

    export interface UploadFileOpts {
        localPath: string
        url: string
    }
    export function uploadFile(o: UploadFileOpts): Promise<void> {
        const data = wx.getFileSystemManager().readFileSync(o.localPath)
        return new Promise((resolve, reject) => {
            wx.request({
                method: 'PUT',
                url: o.url,
                data,
                success: res => {
                    if (res.statusCode >= 400) {
                        reject(res)
                    } else {
                        resolve()
                    }
                },
                fail: reject,
            })
        })
    }
}