export class Problem {
    id: string;
    name: string;

    constructor(id: string, name: string) {
        this.id = id;
        this.name = name;
    }
}


export class LogMessage {
    time: Date;
    source: string;
    text: string;


    constructor(time: Date, source: string, text: string) {
        this.time = time;
        this.source = source;
        this.text = text;
    }
}
