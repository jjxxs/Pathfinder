export class JWebSocket {
    isConnected: boolean = false;
    shutdownRequested: boolean = false;
    private webSocket: WebSocket | null = null;

    constructor(public readonly url: string,
                public readonly timeout: number = 1,
                private onMessage: ((this: JWebSocket, messageEvent: MessageEvent) => any | void) | null,
                private onOpen: ((this: JWebSocket, openEvent: Event) => any | void) | null,
                private onClose: ((this: JWebSocket, closeEvent: CloseEvent) => any | void) | null,
                private onError: ((this: JWebSocket, errorEvent: Event) => any | void) | null) { }

    connect() {
        this.shutdownRequested = false;
        this.webSocket = new WebSocket(this.url);
        this.webSocket.onopen = this._onopen.bind(this);
        this.webSocket.onclose = this._onclose.bind(this);
        this.webSocket.onerror = this._onerror.bind(this);
        this.webSocket.onmessage = this._onmessage.bind(this);
    }

    shutdown() {
        this.shutdownRequested = true;
        if (this.webSocket !== null) {
            this.webSocket.close();
        }
    }

    send(data: any) {
        if (this.webSocket !== null && this.isConnected) {
            this.webSocket.send(data);
        } else {
            console.log("[Warning] WebSocket not connected, dropped message.");
        }
    }

    _onopen(openEvent: Event) {
        this.isConnected = true;
        if (this.onOpen !== null) {
            this.onOpen(openEvent);
        }
    }

    _onmessage(messageEvent: MessageEvent) {
        if (this.onMessage !== null) {
            this.onMessage(messageEvent);
        }
    }

    _onclose(closeEvent: CloseEvent) {
        this.isConnected = false;
        if (this.onClose !== null) {
            this.onClose(closeEvent);
        }

        if (!this.shutdownRequested) {
            setTimeout(this.connect.bind(this), this.timeout);
        }
    }

    _onerror(errorEvent: Event) {
        if (this.onError !== null) {
            this.onError(errorEvent);
        }
    }
}