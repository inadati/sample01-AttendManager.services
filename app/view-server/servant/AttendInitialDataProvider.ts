export module AttendInitialDataProvider {
    export class Summon {
        data: any
        constructor() {
            this.data = {
                isAttend: false,
                inTime: '9:00',
                outTime: '17:00',
            }
        }

        serveData(): any {
            return this.data
        }
    }
}