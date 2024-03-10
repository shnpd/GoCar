export namespace routing {
    export interface drivingOpts {
        trip_id: string
    }
    export function driving(o: drivingOpts) {
        return `/pages/driving/driving?trip_id=${o.trip_id}`
    }
    export interface lockOpts {
        car_id: string
    }
    export function lock(o: lockOpts) {
        return `/pages/lock/lock?car_id=${o.car_id}`
    }
    export interface registerOpts {
        redirect?: string
    }
    export interface registerParams{
        redirectURL:string
    }
    export function register(p?: registerParams) {
        const page = '/pages/register/register'
        if (!p){
            return page
        }
        return `${page}?redirect=${encodeURIComponent(p.redirectURL)}`
    }
    export function mytrips(){
        return '/pages/mytrips/mytrips'
    }
}