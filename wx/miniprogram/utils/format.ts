export function formatDuration(sec: number) {
    const padString = (n: number) => {
        return n < 10 ? '0' + n.toFixed(0) : n.toFixed(0)
    }
    let hour = Math.floor(sec / 3600)
    sec = sec - hour * 3600
    let min = Math.floor(sec / 60)
    sec = sec - min * 60
    return {
        hh: padString(hour),
        mm: padString(min),
        ss: padString(sec),
    }
}
export function formatFee(cents: number): string {
    return (cents / 100).toFixed(2)
}