export namespace models {
	
	export class PomodoroSettings {
	    work_minutes: number;
	    short_break_minutes: number;
	    long_break_minutes: number;
	    long_break_interval: number;
	    auto_start_next: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PomodoroSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.work_minutes = source["work_minutes"];
	        this.short_break_minutes = source["short_break_minutes"];
	        this.long_break_minutes = source["long_break_minutes"];
	        this.long_break_interval = source["long_break_interval"];
	        this.auto_start_next = source["auto_start_next"];
	    }
	}
	export class SessionState {
	    mode: string;
	    // Go type: time
	    current_session_started_at: any;
	    // Go type: time
	    started_at: any;
	    paused: boolean;
	    session_count: number;
	
	    static createFrom(source: any = {}) {
	        return new SessionState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mode = source["mode"];
	        this.current_session_started_at = this.convertValues(source["current_session_started_at"], null);
	        this.started_at = this.convertValues(source["started_at"], null);
	        this.paused = source["paused"];
	        this.session_count = source["session_count"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

