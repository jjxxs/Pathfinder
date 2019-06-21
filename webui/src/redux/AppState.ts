import {Reducer} from "redux";

//**********************************************************
// Actions
//**********************************************************
export enum ActionTypes {
    ON_INIT = "ON_INIT",
    ON_MESSAGE_RECEIVED = "ON_MESSAGE_RECEIVED",
    ON_CONNECT = "ON_CONNECT",
    ON_DISCONNECT = "ON_DISCONNECT",
    ON_POST_IMAGE = "ON_POST_IMAGE",
    ON_POST_COORDINATES = "ON_POST_COORDINATES",
    ON_POST_STATUS = "ON_POST_STATUS",
}

export interface OnInit { type: string }
export const onInit = () => { return {type: ActionTypes.ON_INIT} };

export interface OnConnect { type: string }
export const onConnect = () => { return {type: ActionTypes.ON_CONNECT} };

export interface OnDisconnect { type: string }
export const onDisconnect = () => { return {type: ActionTypes.ON_DISCONNECT} };

export interface OnPostImage { type: string, image: string }
export const onPostImage = (image: string) => { return {type: ActionTypes.ON_POST_IMAGE, image} };

export interface OnPostCoordinates { type: string, coordinates: number[] }
export const onPostCoordinates = (coordinates: number[]) => { return {type: ActionTypes.ON_POST_COORDINATES, coordinates} };

export interface OnPostStatus { type: string, status: Status }
export const onPostStatus = (status: Status) => { return {type: ActionTypes.ON_POST_STATUS, status} };

//**********************************************************
// Messages
//**********************************************************
export enum MessageTypes {
    POST_IMAGE = "PostImage",
    COORDINATES = "Coordinates",
    STATUS = "Status",
}

export interface ImageMessageData {
    image: string;
}

export interface CoordinatesMessageData {
    coordinates: number[];
}

export interface StatusMessageData {
    status: Status;
}

//**********************************************************
// Types
//**********************************************************
export interface AppState {
    connected: boolean;
    image: string;
    points: number[];
    settings: Settings;
    status: Status;
}

export interface Settings {
    server: string;
}

export interface Message {
    type: string;
    data: object;
}

export interface Status {
    algorithm: string;
    problem: string;
    description: string;
    elapsed: string;
    shortest: number;
    running: boolean;
}

//**********************************************************
// Reducer
//**********************************************************
const initialState: AppState = {
    connected: false,
    image: "",
    points: [],
    settings: {server: "ws://localhost:8091/websocket/"},
    status: {algorithm: "", problem: "", description: "", elapsed: "", running: false, shortest: 0},
};

const reducer: Reducer<AppState> = (state = initialState, action) => {
    switch (action.type) {
        case ActionTypes.ON_CONNECT:
            return Object.assign({}, state, {connected: true});
        case ActionTypes.ON_DISCONNECT:
            return initialState;
        case ActionTypes.ON_POST_IMAGE:
            return Object.assign({}, state, {image: (action as OnPostImage).image});
        case ActionTypes.ON_POST_COORDINATES:
            return Object.assign({}, state, {points: (action as OnPostCoordinates).coordinates});
        case ActionTypes.ON_POST_STATUS:
            return Object.assign({}, state, {status: (action as OnPostStatus).status });
        default:
            break;
    }

    return state;
};

export {reducer as appReducer}