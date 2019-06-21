import React from "react";
import {connect} from "react-redux";
import {AppState, Status} from "../redux/AppState";
import Spinner from "./Spinner";

const InfoPanel : React.FunctionComponent<Status> = ({algorithm, problem, description, running, elapsed, shortest}) => {
    let content = <div className={"ml-auto mr-auto"}><Spinner text={""}/></div>;

    // if we haven't received any data yet, show empty
    if (algorithm.length !== 0) {
        content = <div>
            <h4>Algorithm:</h4>
            <h5>{algorithm}</h5>
            <h4>Problem:</h4>
            <h5>{problem}</h5>
            <h4>Description:</h4>
            <h5>{description}</h5>
            <h4>Elapsed:</h4>
            <h5>{elapsed}</h5>
            <h4>Shortest:</h4>
            <h5 className={"pb-0"}>{shortest}</h5>
        </div>;
    }

    return (
        <div className={"info-panel mr-4 mt-auto mb-auto text-right pt-5 pb-5 pr-3"}>
            <div className={"d-flex mt-auto mb-auto"}>
                {content}
            </div>
        </div>);
};

const mapStateToProps = (state: AppState) => {
    return state.status;
};

export default connect(mapStateToProps)(InfoPanel);