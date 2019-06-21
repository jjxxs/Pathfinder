import {JWebSocket} from "../JWebsocket";
import {Dispatch, Middleware, MiddlewareAPI} from "redux";
import {
    ActionTypes,
    AppState, CoordinatesMessageData,
    ImageMessageData,
    Message,
    MessageTypes,
    onConnect,
    onDisconnect, onPostCoordinates,
    onPostImage, onPostStatus, StatusMessageData
} from "./AppState";

let _websocket: JWebSocket;

function onMessageReceived(messageEvent: MessageEvent) {
    return {type: ActionTypes.ON_MESSAGE_RECEIVED, messageEvent}
}

export function serverMiddleware() {
    const webSocketMiddleware: Middleware = ({ getState, dispatch }: MiddlewareAPI) => (next: Dispatch) => action => {
        switch (action.type) {
            case ActionTypes.ON_INIT:
                const url = (getState() as AppState).settings.server;
                console.log("[webSocketMiddleware] ON_INIT, trying to connect to " + url);
                const msgRcvCb = (me: MessageEvent) => { dispatch(onMessageReceived(me)) };
                _websocket = new JWebSocket(url, 1, msgRcvCb, () => dispatch(onConnect()), () => dispatch(onDisconnect()), null);
                _websocket.connect();
                break;
            case ActionTypes.ON_MESSAGE_RECEIVED:
                const msg = JSON.parse(action.messageEvent.data) as Message;
                switch (msg.type) {
                    case MessageTypes.POST_IMAGE:
                        const imageMessageData = msg.data as ImageMessageData;
                        dispatch(onPostImage(imageMessageData.image));
                        break;
                    case MessageTypes.COORDINATES:
                        const coordinatesMessageData = msg.data as CoordinatesMessageData;
                        dispatch(onPostCoordinates(coordinatesMessageData.coordinates));
                        break;
                    case MessageTypes.STATUS:
                        const statusMessageData = msg.data as StatusMessageData;
                        dispatch(onPostStatus(statusMessageData.status));
                        break;
                    default:
                        break;
                }
                break;
        }
        return next(action);
    };

    return webSocketMiddleware
}