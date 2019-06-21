import React from 'react';
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import {connect} from "react-redux";
import {AppState} from "./redux/AppState";
import gopherOffline from "./static/gopher_offline.png";
import Status from "./components/Status";
import CanvasContainer from "./components/Map";

interface PropsFromStore {
    connected: boolean;
}

const ProblemPanel : React.FunctionComponent = () => {
    return <div className={"d-flex w-100"}>
        <Status/>
        <CanvasContainer/>
    </div>;
};

const OfflineMessage : React.FunctionComponent = () => {
    return (
        <div className={"w-100 mt-auto mb-auto text-center"}>
            <img className="offline-banner mb-2" src={gopherOffline} />
            <h4 className={"font-italic"}>Sorry, we're closed!</h4>
        </div>
    );
};

const App : React.FunctionComponent<PropsFromStore> = ({connected}) => {
    const comp = connected ?
        <ProblemPanel/> :
        <OfflineMessage/>;

    return (
      <div className={"container-fluid p-0 content"}>
          <div className={"header"}/>
          <div className={"row justify-content-center"}>
              <div className={"d-flex p-3 app"}>
                  {comp}
              </div>
          </div>
      </div>
    )
};

const mapStateToProps = (state: AppState) => ({connected: state.connected});

export default connect(mapStateToProps)(App);
