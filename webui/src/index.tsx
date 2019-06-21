import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';
import {applyMiddleware, createStore, Store} from "redux";
import {appReducer, AppState, onInit} from "./redux/AppState";
import {serverMiddleware} from "./redux/Middleware";
import {Provider} from "react-redux";

const createMiddleWare = () => applyMiddleware(serverMiddleware());

function configureStore() : Store<AppState> {
    return createStore(appReducer, createMiddleWare());
}

const store = configureStore();
store.dispatch(onInit());
ReactDOM.render(<Provider store={store}><App /></Provider>, document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
